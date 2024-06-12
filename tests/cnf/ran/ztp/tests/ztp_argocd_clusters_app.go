package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/assisted"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/ocm"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranhelper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/helper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/tsparams"
)

var _ = Describe("ZTP Argo CD Clusters Tests", Label(tsparams.LabelArgoCdClustersAppTestCases), func() {
	BeforeEach(func() {
		By("checking that the required clusters are present")
		if !ranhelper.AreClustersPresent([]*clients.Settings{raninittools.HubAPIClient, raninittools.Spoke1APIClient}) {
			Skip("not all of the required clusters are present")
		}
	})

	AfterEach(func() {
		By("resetting the clusters app back to the original settings")
		err := helper.SetGitDetailsInArgoCd(
			tsparams.ArgoCdClustersAppName, tsparams.ArgoCdAppDetails[tsparams.ArgoCdClustersAppName], true, false)
		Expect(err).ToNot(HaveOccurred(), "Failed to reset clusters app git details")
	})

	// 54238 - User modification of klustletaddonconfig via gitops
	It("should override the klusterlet addon configuration and verify the change", reportxml.ID("54238"), func() {
		exists, err := helper.UpdateArgoCdAppGitPath(tsparams.ArgoCdClustersAppName, tsparams.ZtpTestPathClustersApp, true)
		if !exists {
			Skip(err.Error())
		}

		Expect(err).ToNot(HaveOccurred(), "Failed to update Argo CD git path")

		By("validating the klusterlet addon change occurred")
		kac, err := ocm.PullKAC(raninittools.HubAPIClient, ranparam.Spoke1Name, ranparam.Spoke1Name)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull klusterlet addon config")

		err = helper.WaitUntilSearchCollectorEnabled(kac, tsparams.ArgoCdChangeTimeout)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for klusterlet addon config to have search collector enabled")
	})

	// 54238 - User modification of klustletaddonconfig via gitops
	It("should not have nmstateconfig CR when nodeNetwork section does not exist on siteConfig", reportxml.ID("54238"), func() {
		// Update the git path manually so we can potentially skip the test before checking if the NM State
		// Config exists.
		gitDetails := tsparams.ArgoCdAppDetails[tsparams.ArgoCdClustersAppName]
		testGitPath := helper.JoinGitPaths([]string{
			gitDetails.Path,
			tsparams.ZtpTestPathRemoveNmState,
		})

		By("checking if the git path exists")
		if !helper.DoesGitPathExist(gitDetails.Repo, gitDetails.Branch, testGitPath+tsparams.ZtpKustomizationPath) {
			Skip(fmt.Sprintf("git path '%s' could not be found", testGitPath))
		}

		By("checking if the NM state config exists on hub")
		nmStateConfigList, err := assisted.ListNmStateConfigsInAllNamespaces(raninittools.HubAPIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to list NM state configs")
		Expect(nmStateConfigList).ToNot(BeEmpty(), "Failed to find NM state config")

		gitDetails.Path = testGitPath

		By("updating the Argo CD clusters app with the remove NM state git path")
		err = helper.SetGitDetailsInArgoCd(
			tsparams.ArgoCdClustersAppName,
			gitDetails,
			true,
			true)
		Expect(err).ToNot(HaveOccurred(), "Failed to update the Argo CD app with new git details")

		By("validate the NM state config is gone on hub")
		nmStateConfigList, err = assisted.ListNmStateConfigsInAllNamespaces(raninittools.HubAPIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to list NM state configs")
		Expect(nmStateConfigList).To(BeEmpty(), "Found NM state config when it should be gone")
	})
})
