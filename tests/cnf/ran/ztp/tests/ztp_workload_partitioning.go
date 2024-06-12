package tests

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/cgu"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/mco"
	"github.com/openshift-kni/eco-goinfra/pkg/nto" //nolint:mispell
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranhelper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/helper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/tsparams"
	kubeletConfigv1beta1 "k8s.io/kubelet/config/v1beta1"
	"k8s.io/utils/cpuset"
	policiesv1 "open-cluster-management.io/governance-policy-propagator/api/v1"
)

var _ = Describe("ZTP Workload Partitioning Tests", Label(tsparams.LabelWorkloadParitioningTestCases), func() {
	var (
		perfProfile   *nto.Builder
		machineConfig *mco.MCBuilder
		tunedPatch    *nto.TunedBuilder
	)

	BeforeEach(func() {
		By("checking that the required clusters are present")
		if !ranhelper.AreClustersPresent([]*clients.Settings{raninittools.HubAPIClient, raninittools.Spoke1APIClient}) {
			Skip("not all of the required clusters are present")
		}

		By("pulling performance profile if it exists")
		perfProfile, _ = ranhelper.GetPerformanceProfileWithCPUSet(raninittools.Spoke1APIClient, matchesByName)
		if perfProfile == nil {
			Skip(fmt.Sprintf("No performance profile %s found on cluster", tsparams.PerformanceProfileName))
		}

		By("pulling machine config if it exists")
		var err error
		machineConfig, err = mco.PullMachineConfig(raninittools.Spoke1APIClient, tsparams.MachineConfigName)
		if err != nil {
			Skip(fmt.Sprintf("No machine config %s found on cluster", tsparams.MachineConfigName))
		}

		By("pulling tuned patch if it exists")
		tunedPatch, err = nto.PullTuned(raninittools.Spoke1APIClient, tsparams.TunedPatchName, tsparams.TunedNamespace)
		if err != nil {
			Skip(fmt.Sprintf("No tuned %s in namespace %s found on cluster", tsparams.TunedPatchName, tsparams.TunedNamespace))
		}
	})

	AfterEach(func() {
		By("resetting the policies app to the original settings")
		err := helper.SetGitDetailsInArgoCd(
			tsparams.ArgoCdPoliciesAppName, tsparams.ArgoCdAppDetails[tsparams.ArgoCdPoliciesAppName], true, false)
		Expect(err).ToNot(HaveOccurred(), "Failed to reset policies app git details")

		By("removing the cgu if it exists")
		updateCgu, err := cgu.Pull(raninittools.HubAPIClient, tsparams.WorkloadPartitioningCguName, tsparams.TestNamespace)
		if err == nil {
			_, err = updateCgu.DeleteAndWait(5 * time.Minute)
			Expect(err).ToNot(HaveOccurred(), "Failed to delete update CGU")
		}

		By("restoring original config for workload partitioning")
		if machineConfig == nil || perfProfile == nil || tunedPatch == nil {
			Skip("No need to restore workload parititoning")
		}

		By("pulling current workload partitioning resources")
		perfProfileTemp, err := ranhelper.GetPerformanceProfileWithCPUSet(raninittools.Spoke1APIClient, matchesByName)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull performance profile")

		tunedPatchTemp, err := nto.PullTuned(raninittools.Spoke1APIClient, tsparams.TunedPatchName, tsparams.TunedNamespace)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull tuned patch")

		machineConfigTemp, err := mco.PullMachineConfig(raninittools.Spoke1APIClient, tsparams.MachineConfigName)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull machine config")

		By("updating workload partitioning resources")
		perfProfileTemp.Definition = perfProfile.Definition
		_, err = perfProfileTemp.Update(true)
		Expect(err).ToNot(HaveOccurred(), "Failed to update performance profile")

		tunedPatchTemp.Definition.Spec = tunedPatch.Definition.Spec

		if tunedPatch.Definition.GetAnnotations() != nil {
			tunedPatchTemp.Definition.SetAnnotations(tunedPatchTemp.Definition.Annotations)
		}

		if tunedPatch.Definition.GetLabels() != nil {
			tunedPatchTemp.Definition.SetLabels(tunedPatchTemp.Definition.Labels)
		}

		_, err = tunedPatchTemp.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update tuned patch")

		machineConfigTemp.Definition.Spec = machineConfig.Definition.Spec

		if machineConfig.Definition.GetAnnotations() != nil {
			machineConfigTemp.Definition.SetAnnotations(machineConfig.Definition.Annotations)
		}

		if machineConfig.Definition.GetLabels() != nil {
			machineConfigTemp.Definition.SetLabels(machineConfigTemp.Definition.Labels)
		}

		_, err = machineConfigTemp.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update machine config")

		By("waiting for SNO to be functional after rebooting from changes")
		err = helper.WaitForNodeFunctional(raninittools.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for SNO to be functional again")

		validateWorkloadPartition()
	})

	// 54235 - Supporting changes to workload partitioning cpuset at day-n.
	It("should update cpu partitioning with the new reserved cpus", reportxml.ID("54235"), func() {
		By("updating Argo CD policies app")
		exists, err := helper.UpdateArgoCdAppGitPath(
			tsparams.ArgoCdPoliciesAppName, tsparams.ZtpTestPathWorkloadPartitioning, true)
		if !exists {
			Skip(err.Error())
		}

		Expect(err).ToNot(HaveOccurred(), "Failed to update Argo CD git path")

		By("waiting for the policy to exist")
		policy, err := helper.WaitForPolicyToExist(
			raninittools.HubAPIClient,
			tsparams.WorkloadPartitioningPolicyName,
			tsparams.TestNamespace,
			tsparams.ArgoCdChangeTimeout)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for policy to exist")

		By("waiting for the policy to be NonCompliant")
		err = policy.WaitUntilComplianceState(policiesv1.NonCompliant, tsparams.ArgoCdChangeTimeout)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for policy to be Compliant")

		By("creating CGU to apply changes")
		cguBuilder := cgu.NewCguBuilder(
			raninittools.HubAPIClient, tsparams.WorkloadPartitioningCguName, tsparams.TestNamespace, 1).
			WithCluster(ranparam.Spoke1Name).
			WithManagedPolicy(tsparams.WorkloadPartitioningPolicyName)
		cguBuilder.Definition.Spec.RemediationStrategy.Timeout = 10

		_, err = cguBuilder.Create()
		Expect(err).ToNot(HaveOccurred(), "Failed to create CGU")

		By("waiting for SNO to be functional after rebooting from changes")
		err = helper.WaitForNodeFunctional(raninittools.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for SNO to be functional again")

		validateWorkloadPartition()
	})
})

func matchesByName(profile *nto.Builder) bool {
	return profile.Definition.Name == tsparams.PerformanceProfileName
}

// validateWorkloadPartition ensures that the workload partitioning is setup properly and processes are affined to the
// right cpus.
func validateWorkloadPartition() {
	perfProfile, err := ranhelper.GetPerformanceProfileWithCPUSet(raninittools.Spoke1APIClient, matchesByName)
	Expect(err).ToNot(HaveOccurred(), "Failed to pull performance profile")

	reservedCPUSet, err := cpuset.Parse(string(*perfProfile.Object.Spec.CPU.Reserved))
	Expect(err).ToNot(HaveOccurred(), "Failed to parse reserved CPU set")

	By("checking kubeletconfig reservedSystemCPUs")

	kubeletConfig, err := mco.PullKubeletConfig(raninittools.Spoke1APIClient, tsparams.KubeletConfigName)
	Expect(err).ToNot(HaveOccurred(), "Failed to pull kubelet config")

	kubeletConfiguration, err := ranhelper.UnmarshalRaw[kubeletConfigv1beta1.KubeletConfiguration](
		kubeletConfig.Object.Spec.KubeletConfig.Raw)
	Expect(err).ToNot(HaveOccurred(), "Failed to unmarshal kubelet configuration")
	Expect(kubeletConfiguration.ReservedSystemCPUs).
		To(Equal(reservedCPUSet.String()), "KubeletConfig did not match reservedSystemCPUs")

	By("checking new cpu processes affinities")

	processNames := []string{"crio", "kubelet", "ovn"}
	err = helper.CheckAffinitiesByProcessMatch(raninittools.Spoke1APIClient, processNames, reservedCPUSet)
	Expect(err).ToNot(HaveOccurred(), "Failed to check cpu processes affinities")

	By("checking cpuset for all running containers via crictl inspect")

	containersInfo, err := helper.GetContainersInfo(raninittools.Spoke1APIClient)
	Expect(err).ToNot(HaveOccurred(), "Failed to get info for running containers")

	err = helper.CheckPodsAffinity(helper.GetManagementContainersInfo(containersInfo), reservedCPUSet)
	Expect(err).ToNot(HaveOccurred(), "Failed to check pod affinities for management containers")

	By("checking os daemon is pinned to reserved cpus")

	_, err = helper.CheckCPUAffinityOnNonKernelPids(raninittools.Spoke1APIClient, reservedCPUSet)
	Expect(err).ToNot(HaveOccurred(), "Processes not pinned to reserved cpus")
}
