package policy

import (
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/namespace"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	. "github.com/openshift-kni/eco-gotests/tests/cnf/core/network/internal/netinittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/core/network/policy/internal/tsparams"
	_ "github.com/openshift-kni/eco-gotests/tests/cnf/core/network/policy/tests"
	"github.com/openshift-kni/eco-gotests/tests/internal/cluster"
	"github.com/openshift-kni/eco-gotests/tests/internal/params"
	"github.com/openshift-kni/eco-gotests/tests/internal/reporter"
)

var (
	_, currentFile, _, _ = runtime.Caller(0)
	testNS               = namespace.NewBuilder(APIClient, tsparams.TestNamespaceName)
)

func TestLB(t *testing.T) {
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = NetConfig.GetJunitReportPath(currentFile)

	RegisterFailHandler(Fail)
	RunSpecs(t, "NetworkPolicy", Label(tsparams.Labels...), reporterConfig)
}

var _ = BeforeSuite(func() {
	By("Creating test namespace with privileged labels")
	for key, value := range params.PrivilegedNSLabels {
		testNS.WithLabel(key, value)
	}
	_, err := testNS.Create()
	Expect(err).ToNot(HaveOccurred(), "error to create test namespace")
	By("Pulling test images on cluster before running test cases")
	err = cluster.PullTestImageOnNodes(APIClient, NetConfig.WorkerLabel, NetConfig.CnfNetTestContainer, 300)
	Expect(err).ToNot(HaveOccurred(), "Failed to pull test image on nodes")
})

var _ = AfterSuite(func() {
	By("Deleting test namespace")
	err := testNS.DeleteAndWait(tsparams.DefaultTimeout)
	Expect(err).ToNot(HaveOccurred(), "error to delete test namespace")
})

var _ = JustAfterEach(func() {
	reporter.ReportIfFailed(
		CurrentSpecReport(), currentFile, tsparams.ReporterNamespacesToDump, tsparams.ReporterCRDsToDump, clients.SetScheme)
})

var _ = ReportAfterSuite("", func(report Report) {
	reportxml.Create(report, NetConfig.GetReportPath(), NetConfig.TCPrefix)
})
