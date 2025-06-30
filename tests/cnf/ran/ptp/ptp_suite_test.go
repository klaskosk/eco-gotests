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
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ptp/internal/consumer"
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

	By("deploying consumers")
	err := consumer.DeployConsumersOnWorkers(RANConfig.Spoke1APIClient)
	Expect(err).ToNot(HaveOccurred(), "Failed to deploy consumers on workers")

	// GetNicDriver gets the driver for a given interface by running cmd via a PTP pod.
	// func GetNicDriver(ptpPod corev1.Pod, ifName string) (string, error) {
	// 	cmd := fmt.Sprintf("ethtool -i %s | grep --color=no driver | awk '{print $2}'", ifName)
	// 	out, err := pod.ExecCommand(helper.Apiclient, ptpPod, []string{"/bin/bash", "-c", cmd}, parameters.PtpContainerName)

	// 	if nil != err {
	// 		return "", err
	// 	}

	// 	return strings.Trim(out.String(), "\n"), nil
	// }
	By("increasing thresholds to 200 on mlx")
})

var _ = AfterSuite(func() {
	By("restoring PTP configs")

	By("removing consumers")
	err := consumer.CleanupConsumersOnWorkers(RANConfig.Spoke1APIClient)
	Expect(err).ToNot(HaveOccurred(), "Failed to cleanup consumers on workers")
})

var _ = JustAfterEach(func() {
	reporter.ReportIfFailed(
		CurrentSpecReport(), currentFile, tsparams.ReporterSpokeNamespacesToDump, tsparams.ReporterSpokeCRsToDump)
})

var _ = ReportAfterSuite("", func(report Report) {
	reportxml.Create(report, RANConfig.GetReportPath(), RANConfig.TCPrefix)
})
