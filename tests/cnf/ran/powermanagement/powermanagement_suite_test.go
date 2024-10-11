package powermanagement

import (
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/namespace"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	. "github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/powermanagement/internal/tsparams"
	_ "github.com/openshift-kni/eco-gotests/tests/cnf/ran/powermanagement/tests"
	"github.com/openshift-kni/eco-gotests/tests/internal/reporter"
)

var _, currentFile, _, _ = runtime.Caller(0)

func TestPowerSave(t *testing.T) {
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = RANConfig.GetJunitReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Power Management Test Suite", Label(tsparams.Labels...), reporterConfig)
}

var _ = BeforeSuite(func() {
	// Cleanup and create test namespace
	testNamespace := namespace.NewBuilder(Spoke1APIClient, tsparams.TestingNamespace).
		WithLabel("pod-security.kubernetes.io/enforce", "baseline")

	By("deleting and recreating test namespace to ensure blank slate")
	err := testNamespace.DeleteAndWait(tsparams.PowerSaveTimeout)
	Expect(err).ToNot(HaveOccurred(), "Failed to delete namespace ", tsparams.TestingNamespace)

	_, err = testNamespace.Create()
	Expect(err).ToNot(HaveOccurred(), "Failed to create namespace ", tsparams.TestingNamespace)
})

var _ = AfterSuite(func() {
	By("deleting test namespace to clean up test suite")
	testNamespace := namespace.NewBuilder(Spoke1APIClient, tsparams.TestingNamespace)
	err := testNamespace.DeleteAndWait(tsparams.PowerSaveTimeout)
	Expect(err).ToNot(HaveOccurred(), "Failed to delete namespace ", tsparams.TestingNamespace)

})

var _ = JustAfterEach(func() {
	reporter.ReportIfFailed(
		CurrentSpecReport(), currentFile, tsparams.ReporterNamespacesToDump, tsparams.ReporterCRsToDump)
})

var _ = ReportAfterSuite("", func(report Report) {
	reportxml.Create(report, RANConfig.GetReportPath(), RANConfig.TCPrefix)
})
