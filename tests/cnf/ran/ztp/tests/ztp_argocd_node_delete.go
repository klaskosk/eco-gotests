package tests

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/bmh"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/mco"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranhelper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/helper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/tsparams"
)

var _ = Describe("ZTP Argo CD Node Deletion Tests", Label(tsparams.LabelArgoCdNodeDeletionTestCases), func() {
	var (
		plusOneNodeName   string
		bmhNamespace      string
		earlyReturnNonSNO = true
	)

	BeforeEach(func() {
		By("checking that the required clusters are present")
		if !ranhelper.AreClustersPresent([]*clients.Settings{raninittools.HubAPIClient, raninittools.Spoke1APIClient}) {
			Skip("not all of the required clusters are present")
		}

		By("checking the ZTP version")
		versionInRange, err := ranhelper.IsVersionStringInRange(ranparam.ZtpVersion, "4.14", "")
		Expect(err).ToNot(HaveOccurred(), "Failed to check if ZTP version is in range")

		if !versionInRange {
			Skip("ZTP node deletion tests require ZTP version of at least 4.14")
		}

		By("checking that the 'worker' mcp is ready")
		mcp, err := mco.Pull(raninittools.Spoke1APIClient, "worker")
		Expect(err).ToNot(HaveOccurred(), "Failed to pull 'worker' MCP")

		if mcp.Definition.Status.ReadyMachineCount <= 0 {
			Skip("ZTP node deletion tests require ready 'worker' MCP")
		}

		By("checking that the cluster contains a master/control-plane and worker node")
		snoPlusOne, err := helper.IsSnoPlusOne(raninittools.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to check if cluster is SNO+1")

		if !snoPlusOne {
			Skip("Cluster does not contain a single master and a single worker node")
		}

		earlyReturnNonSNO = false

		plusOneNodeName, err = helper.GetPlusOneWorkerName(raninittools.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to get SNO+1 worker name")

		bmhNamespace, err = helper.GetBmhNamespace(raninittools.HubAPIClient, plusOneNodeName)
		Expect(err).ToNot(HaveOccurred(), "Failed to get BMH namespace")
		Expect(bmhNamespace).ToNot(BeEmpty(), "BMH namespace cannot be empty")
	})

	AfterEach(func() {
		if earlyReturnNonSNO {
			return
		}

		By("resetting the clusters app back to the original settings")
		err := helper.SetGitDetailsInArgoCd(
			tsparams.ArgoCdClustersAppName, tsparams.ArgoCdAppDetails[tsparams.ArgoCdClustersAppName], true, true)
		Expect(err).ToNot(HaveOccurred(), "Failed to reset clusters app git details")

		By("checking that the cluster is back to SNO+1")
		err = helper.WaitForNumberOfNodes(raninittools.Spoke1APIClient, 2, 45*time.Minute)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for cluster to return to 2 nodes")

		snoPlusOne, err := helper.IsSnoPlusOne(raninittools.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to check if cluster is SNO+1")
		Expect(snoPlusOne).To(BeTrue(), "Cluster is no longer SNO+1")
	})

	// 72463 - Delete and re-add a worker node from cluster
	It("should delete a worker node from the cluster", reportxml.ID("72463"), func() {
		By("updating the Argo CD git path to apply crAnnotation")
		exists, err := helper.UpdateArgoCdAppGitPath(
			tsparams.ArgoCdClustersAppName, tsparams.ZtpTestPathNodeDeleteAddAnnotation, true)
		if !exists {
			Skip(err.Error())
		}

		Expect(err).ToNot(HaveOccurred(), "Failed to update Argo CD git path")

		By("waiting for the crAnnotation to be added to the worker node")
		bareMetalHost, err := bmh.Pull(raninittools.HubAPIClient, plusOneNodeName, bmhNamespace)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull BMH")

		err = helper.WaitForBMHAnnotation(bareMetalHost, tsparams.NodeDeletionCrAnnotation, tsparams.ArgoCdChangeTimeout)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for BMH annotation")

		By("updating the Argo CD app to apply the suppression to the spec")
		exists, err = helper.UpdateArgoCdAppGitPath(
			tsparams.ArgoCdClustersAppName, tsparams.ZtpTestPathNodeDeleteAddSuppression, false)
		if !exists {
			Skip(err.Error())
		}

		Expect(err).ToNot(HaveOccurred(), "Failed to update Argo CD git path")

		By("waiting for the worker node to be removed")
		err = helper.WaitForBMHDeprovisioning(raninittools.HubAPIClient, plusOneNodeName, bmhNamespace, 60*time.Minute)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for worker BMH to be deprovisioned")

		By("checking that the cluster is healthy")
		healthy, err := helper.IsClusterStable(raninittools.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to check if spoke cluster is healthy")
		Expect(healthy).To(BeTrue(), "Spoke cluster was not healthy")
	})
})
