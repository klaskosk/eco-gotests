package ptp

import (
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/rancluster"
	. "github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
	_ "github.com/openshift-kni/eco-gotests/tests/cnf/ran/ptp/tests"
	"github.com/openshift-kni/eco-gotests/tests/internal/reporter"
)

var _, currentFile, _, _ = runtime.Caller(0)

func TestPTP(t *testing.T) {
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = RANConfig.GetJunitReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "RAN PTP Suite", Label(tsparams.Labels...), reporterConfig)
}

var _ = BeforeSuite(func() {
	By("checking that the spoke 1 cluster is present")
	isSpoke1Present := rancluster.AreClustersPresent([]*clients.Settings{Spoke1APIClient})
	Expect(isSpoke1Present).To(BeTrue(), "Spoke 1 cluster must be present for PTP tests")
})

var _ = JustAfterEach(func() {
	reporter.ReportIfFailed(
		CurrentSpecReport(), currentFile, tsparams.ReporterSpokeNamespacesToDump, tsparams.ReporterSpokeCRsToDump)
})

var _ = ReportAfterSuite("", func(report Report) {
	reportxml.Create(report, RANConfig.GetReportPath(), RANConfig.TCPrefix)
})
