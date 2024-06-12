package tests

import (
	"strings"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/cgu"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/namespace"
	"github.com/openshift-kni/eco-goinfra/pkg/ocm"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-goinfra/pkg/secret"
	"github.com/openshift-kni/eco-goinfra/pkg/sriov"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranhelper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/helper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/tsparams"
	corev1 "k8s.io/api/core/v1"
	policiesv1 "open-cluster-management.io/governance-policy-propagator/api/v1"
)

var _ = Describe("ZTP Argo CD Hub Templating Tests", Label(tsparams.LabelArgoCdHubTemplatingTestCases), func() {
	BeforeEach(func() {
		By("checking that the required clusters are present")
		if !ranhelper.AreClustersPresent([]*clients.Settings{raninittools.HubAPIClient, raninittools.Spoke1APIClient}) {
			Skip("not all of the required clusters are present")
		}
	})

	AfterEach(func() {
		By("resetting the policies app back to the original settings")
		err := helper.SetGitDetailsInArgoCd(
			tsparams.ArgoCdPoliciesAppName, tsparams.ArgoCdAppDetails[tsparams.ArgoCdPoliciesAppName], true, false)
		Expect(err).ToNot(HaveOccurred(), "Failed to reset the git details for the policies app")

		By("removing the hub templating leftovers if any exist")
		network, err := sriov.PullNetwork(raninittools.Spoke1APIClient, tsparams.TestNamespace, tsparams.SrIovNetworkOperator)
		if err == nil {
			err = network.DeleteAndWait(tsparams.ArgoCdChangeTimeout)
			Expect(err).ToNot(HaveOccurred(), "Failed to delete SR-IOV network")
		}

		By("removing the CGU if it exists")
		cguBuilder, err := cgu.Pull(
			raninittools.HubAPIClient, tsparams.HubTemplatingCguName, tsparams.HubTemplatingCguNamespace)
		if err == nil {
			_, err = cguBuilder.DeleteAndWait(5 * time.Minute)
			Expect(err).ToNot(HaveOccurred(), "Failed to delete and wait for CGU to be deleted")
		}

		By("removing test secret if it exists")
		secretBuilder, err := secret.Pull(raninittools.HubAPIClient, ranparam.Spoke1Name+"-sriovdata", tsparams.TestNamespace)
		if err == nil {
			err = secretBuilder.Delete()
			Expect(err).ToNot(HaveOccurred(), "Failed to delete test secret")
		}
	})

	When("TALM version is at most 4.15", func() {
		BeforeEach(func() {
			By("checking the ZTP version")
			versionInRange, err := ranhelper.IsVersionStringInRange(ranparam.ZtpVersion, "", "4.15")
			Expect(err).ToNot(HaveOccurred(), "Failed to check if ZTP version is in range")

			if !versionInRange {
				Skip("This test requires a ZTP version of at least 4.14")
			}
		})

		// 54240 - Hub-side ACM templating with TALM
		It("should report an error for using printf function where not allowed", reportxml.ID("54240"), func() {
			setupHubTemplateTest(tsparams.ZtpTestPathTemplatingPrintf)

			By("validating TALM reported a policy error")
			assertTalmPodLog(raninittools.HubAPIClient, "printf variable is not supported in the template function Name field")
		})

		// 54240 - Hub-side ACM templating with TALM
		It("should report an error for using fromsecret function where not allowed", reportxml.ID("54240"), func() {
			// Since we need to create the secret first, we need to manually create the namespace instead of relying
			// on the site config.
			_, err := namespace.NewBuilder(raninittools.HubAPIClient, tsparams.TestNamespace).Create()
			Expect(err).ToNot(HaveOccurred(), "Failed to create test namespace")

			secretBuilder := secret.NewBuilder(
				raninittools.HubAPIClient, ranparam.Spoke1Name+"-sriovdata", tsparams.TestNamespace, corev1.SecretTypeOpaque)
			secretBuilder.Definition.StringData = map[string]string{"vlan": "MTEwCg=="}
			_, err = secretBuilder.Create()
			Expect(err).ToNot(HaveOccurred(), "Failed to create secret")

			setupHubTemplateTest(tsparams.ZtpTestPathTemplatingFromSecret)

			By("validating TALM reported a policy error")
			assertTalmPodLog(raninittools.HubAPIClient, "template function is not supported in TALM")
		})

		// 54240 - Hub-side ACM templating with TALM
		It("should report an error for using autoindent function where not allowed", reportxml.ID("54240"), func() {
			setupHubTemplateTest(tsparams.ZtpTestPathTemplatingAutoIndent)

			By("validating TALM reported a policy error")
			assertTalmPodLog(raninittools.HubAPIClient, "policy has hub template error")

			By("validating the specific error using the policy message")
			err := helper.WaitForPolicyMessageToContainSubstring(
				raninittools.HubAPIClient,
				tsparams.TestNamespace+"."+tsparams.HubTemplatingPolicyName,
				ranparam.Spoke1Name,
				"wrong type for value; expected string; got int")
			Expect(err).ToNot(HaveOccurred(), "Failed to validate error using policy message")
		})

		It("should report an error for using invalid lookup hub template function", func() {
			setupHubTemplateTest(tsparams.ZtpTestPathTemplatingLookupInvalid)

			By("validating CGU reported an invalid policy error")
			cguBuilder, err := cgu.Pull(
				raninittools.HubAPIClient, tsparams.HubTemplatingCguName, tsparams.HubTemplatingCguNamespace)
			Expect(err).ToNot(HaveOccurred(), "Failed to pull hub templating CGU")

			_, err = cguBuilder.WaitForCondition(tsparams.InvalidManagedPoliciesCondition, time.Minute)
			Expect(err).ToNot(HaveOccurred(), "Failed to wait for CGU to report invalid policy error")

			By("validating TALM reported a policy error")
			assertTalmPodLog(raninittools.HubAPIClient, "template function only allows the resource with apiVersion in"+
				" 'cluster.open-cluster-management.io', kind 'ManagedCluster' and empty namespace")
		})
	})

	// 54240 - Hub-side ACM templating with TALM
	It("should create the policy successfully with a valid template", reportxml.ID("54240"), func() {
		setupHubTemplateTest(tsparams.ZtpTestPathTemplatingValid)

		By("validate the policy reaches compliant status")
		policy, err := ocm.PullPolicy(raninittools.HubAPIClient, tsparams.HubTemplatingPolicyName, tsparams.TestNamespace)
		Expect(err).ToNot(HaveOccurred(), "Failed to get policy from hub cluster")

		err = policy.WaitUntilComplianceState(policiesv1.Compliant, tsparams.ArgoCdChangeTimeout)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for policy to become compliant")
	})
})

// setupHubTemplateTest extracts the core setup logic for the hub templating test cases.
func setupHubTemplateTest(ztpTestPath string) {
	By("updating the Argo CD git path")

	exists, err := helper.UpdateArgoCdAppGitPath(tsparams.ArgoCdPoliciesAppName, ztpTestPath, true)
	if !exists {
		Skip(err.Error())
	}

	Expect(err).ToNot(HaveOccurred(), "Failed to update Argo CD git path")

	By("waiting for the policy to exist")

	_, err = helper.WaitForPolicyToExist(
		raninittools.HubAPIClient, tsparams.HubTemplatingPolicyName, tsparams.TestNamespace, tsparams.ArgoCdChangeTimeout)
	Expect(err).ToNot(HaveOccurred(), "Failed to wait for hub templating policy to be created")

	By("creating the CGU")

	cguBuilder := cgu.NewCguBuilder(
		raninittools.HubAPIClient, tsparams.HubTemplatingCguName, tsparams.HubTemplatingCguNamespace, 1).
		WithCluster(ranparam.Spoke1Name).
		WithManagedPolicy(tsparams.HubTemplatingPolicyName)
	cguBuilder.Definition.Spec.RemediationStrategy.Timeout = 10

	_, err = cguBuilder.Create()
	Expect(err).ToNot(HaveOccurred(), "Failed to create hub templating CGU")
}

// assertTalmPodLog asserts that the TALM pod log contains the expected substring.
func assertTalmPodLog(client *clients.Settings, expectedSubstring string) {
	glog.V(tsparams.LogLevel).Infof("Waiting for TALM log to report: '%s'", expectedSubstring)

	Eventually(func() string {
		podList, err := pod.List(client, ranparam.OpenshiftOperatorNamespace)
		Expect(err).ToNot(HaveOccurred(), "Failed to list pods is openshift operator namespace")
		Expect(podList).ToNot(BeEmpty(), "Failed to find any pods in the openshift operator namespace")

		var podLog string

		for _, podBuilder := range podList {
			if strings.HasPrefix(podBuilder.Definition.GetName(), tsparams.TalmHubPodName) {
				podLog, err = podBuilder.GetLog(1*time.Minute, ranparam.TalmContainerName)
				Expect(err).ToNot(HaveOccurred(), "Failed to get TALM pod log")

				break
			}
		}

		return podLog
	}, tsparams.ArgoCdChangeTimeout, tsparams.ArgoCdChangeInterval).
		Should(ContainSubstring(expectedSubstring), "Failed to assert TALM pod log contains %s", expectedSubstring)
}
