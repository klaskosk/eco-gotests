package tests

import (
	"slices"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/argocd"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/deployment"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranhelper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/helper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/tsparams"
	corev1 "k8s.io/api/core/v1"
	policiesv1 "open-cluster-management.io/governance-policy-propagator/api/v1"
)

var _ = Describe("ZTP Argo CD ACM CR Tests", Label(tsparams.LabelArgoCdAcmCrsTestCases), func() {
	var (
		acmPolicyGeneratorImage        string
		oldAcmPolicyGeneratorContainer corev1.Container
	)

	BeforeEach(func() {
		By("checking that the required clusters are present")
		if !ranhelper.AreClustersPresent([]*clients.Settings{raninittools.HubAPIClient, raninittools.Spoke1APIClient}) {
			Skip("not all of the required clusters are present")
		}

		By("determining the container image for ACM CR integration")
		multiClusterDeployment, err := deployment.Pull(
			raninittools.HubAPIClient, tsparams.MultiClusterHubOperator, ranparam.AcmOperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Failed to get multi cluster operator deployment")

		acmPolicyGeneratorImage = getContainerImageFromDeploymentEnvironment(
			multiClusterDeployment, tsparams.MultiClusterHubOperator, "OPERAND_IMAGE_MULTICLUSTER_OPERATORS_SUBSCRIPTION")
		Expect(acmPolicyGeneratorImage).ToNot(BeEmpty(), "Failed to find ACM policy generator container image")

		glog.V(tsparams.LogLevel).Infof("Found ACM policy generator container image: '%s'", acmPolicyGeneratorImage)

		By("updating Argo CD to allow ACM CRs")
		argoCd, err := argocd.Pull(raninittools.HubAPIClient, ranparam.OpenshiftGitops, ranparam.OpenshiftGitops)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull Argo CD instance")

		// The default container here is based off the documentation:
		//nolint:lll
		// https://github.com/openshift-kni/cnf-features-deploy/blob/master/ztp/gitops-subscriptions/ACMPolicyGeneratorIntergration.md#openshift-gitopsargocd

		acmPolicyGeneratorContainer, err := pod.NewContainerBuilder(
			tsparams.AcmPolicyGeneratorName,
			acmPolicyGeneratorImage,
			[]string{"/bin/bash", "-c", "cp -r /etc/kustomize /.config"}).
			WithImagePullPolicy(corev1.PullAlways).
			WithVolumeMount(corev1.VolumeMount{Name: "kustomize", MountPath: "/.config"}).
			GetContainerCfg()
		Expect(err).ToNot(HaveOccurred(), "Failed to get default container config")

		for index, container := range argoCd.Definition.Spec.Repo.InitContainers {
			if container.Name == tsparams.AcmPolicyGeneratorName {
				oldAcmPolicyGeneratorContainer = container

				acmPolicyGeneratorContainer = container.DeepCopy()
				acmPolicyGeneratorContainer.Image = acmPolicyGeneratorImage

				argoCd.Definition.Spec.Repo.InitContainers = slices.Delete(
					argoCd.Definition.Spec.Repo.InitContainers, index, index+1)

				break
			}
		}

		// Move the ACM container to be first.
		argoCd.Definition.Spec.Repo.InitContainers = append(
			[]corev1.Container{*acmPolicyGeneratorContainer}, argoCd.Definition.Spec.Repo.InitContainers...)

		_, err = argoCd.Update(true)
		Expect(err).ToNot(HaveOccurred(), "Failed to update the Argo CD instance")

		By("Waiting a minute for Argo CD to catch up")
		time.Sleep(1 * time.Minute)
	})

	AfterEach(func() {
		By("resetting the policies app back to the original settings")
		err := helper.SetGitDetailsInArgoCd(
			tsparams.ArgoCdPoliciesAppName, tsparams.ArgoCdAppDetails[tsparams.ArgoCdPoliciesAppName], true, false)
		Expect(err).ToNot(HaveOccurred(), "Failed to reset the git details for the policies app")

		By("reverting Argo CD patch for ACM CRs")
		argoCd, err := argocd.Pull(raninittools.HubAPIClient, ranparam.OpenshiftGitops, ranparam.OpenshiftGitops)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull Argo CD instance")

		for i, container := range argoCd.Definition.Spec.Repo.InitContainers {
			if container.Name == tsparams.AcmPolicyGeneratorName {
				argoCd.Definition.Spec.Repo.InitContainers = slices.Delete(argoCd.Definition.Spec.Repo.InitContainers, i, i+1)

				if oldAcmPolicyGeneratorContainer.Name != "" {
					argoCd.Definition.Spec.Repo.InitContainers = append(
						argoCd.Definition.Spec.Repo.InitContainers, oldAcmPolicyGeneratorContainer)
				}

				_, err = argoCd.Update(true)
				Expect(err).ToNot(HaveOccurred(), "Failed to update Argo CD instance")

				By("Waiting a minute for Argo CD to catch up")
				time.Sleep(1 * time.Minute)

				break
			}
		}
	})

	// 54236 - Evaluating use of ACM's version of PolicyGenTemplates with our ZTP flow. This enables user created
	// content that does not depend on our ZTP container but works "seamlessly" with it.
	It("should use ACM CRs to template a policy, deploy the policy, and validate it was successful",
		reportxml.ID("54236"), func() {
			exists, err := helper.UpdateArgoCdAppGitPath(tsparams.ArgoCdPoliciesAppName, tsparams.ZtpTestPathAcmCrs, true)
			if !exists {
				Skip(err.Error())
			}

			Expect(err).ToNot(HaveOccurred(), "Failed to update Argo CD git path")

			By("waiting for policies to be created")
			policy, err := helper.WaitForPolicyToExist(
				raninittools.HubAPIClient, tsparams.AcmCrsPolicyName, tsparams.TestNamespace, tsparams.ArgoCdChangeTimeout)
			Expect(err).ToNot(HaveOccurred(), "Failed to wait for the ACM CRs policy to be created")

			By("validating the policy was created and wait for it to finish")
			err = policy.WaitUntilComplianceState(policiesv1.NonCompliant, 1*time.Minute)
			Expect(err).ToNot(HaveOccurred(), "Failed to wait for ACM CRs policy to be non-compliant")
		})
})

// getContainerImageFromDeploymentEnvironment gets the value of an environment variable from a specific container in a
// deployment.
func getContainerImageFromDeploymentEnvironment(
	deploymentBuilder *deployment.Builder, containerName, envName string) string {
	for _, container := range deploymentBuilder.Definition.Spec.Template.Spec.Containers {
		if container.Name == containerName {
			for _, envVar := range container.Env {
				if envVar.Name == envName {
					return envVar.Value
				}
			}
		}
	}

	return ""
}
