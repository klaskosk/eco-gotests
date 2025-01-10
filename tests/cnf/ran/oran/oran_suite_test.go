package oran

import (
	"fmt"
	"path"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	_ "github.com/openshift-kni/eco-gotests/tests/cnf/ran/gitopsztp/tests"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/rancluster"
	. "github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"github.com/openshift-kni/eco-gotests/tests/internal/reporter"
)

var _, currentFile, _, _ = runtime.Caller(0)

func TestORAN(t *testing.T) {
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = RANConfig.GetJunitReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "RAN O-RAN Suite", Label(tsparams.Labels...), reporterConfig)
}

var _ = BeforeSuite(func() {
	By("checking that the hub cluster is present")
	if !rancluster.AreClustersPresent([]*clients.Settings{HubAPIClient}) {
		Skip("the hub cluster is not present")
	}
})

var _ = JustAfterEach(func() {
	var (
		currentDir, currentFilename = path.Split(currentFile)
		hubReportPath               = fmt.Sprintf("%shub_%s", currentDir, currentFilename)
		report                      = CurrentSpecReport()
	)

	if Spoke1APIClient != nil {
		reporter.ReportIfFailed(
			report, currentFile, tsparams.ReporterSpokeNamespacesToDump, tsparams.ReporterSpokeCRsToDump)
	}

	reporter.ReportIfFailedOnCluster(
		RANConfig.HubKubeconfig,
		report,
		hubReportPath,
		tsparams.ReporterHubNamespacesToDump,
		tsparams.ReporterHubCRsToDump)
})

var _ = ReportAfterSuite("", func(report Report) {
	reportxml.Create(report, RANConfig.GetReportPath(), RANConfig.TCPrefix)
})
