package tests

import (
	"strings"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/mco"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranhelper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/tsparams"
)

var _ = Describe("ZTP Machine Config Tests", Label(tsparams.LabelMachineConfigTestCases), func() {
	BeforeEach(func() {
		By("checking that the required clusters are present")
		if !ranhelper.AreClustersPresent([]*clients.Settings{raninittools.HubAPIClient, raninittools.Spoke1APIClient}) {
			Skip("not all of the required clusters are present")
		}
	})

	// 54239 - Annotation on generated CRs for traceability
	It("should find the ztp annotation present in the machine configs", reportxml.ID("54239"), func() {
		machineConfigsToCheck := []string{
			"02-master-workload-partitioning",
			"container-mount-namespace-and-kubelet-conf-master",
			"container-mount-namespace-and-kubelet-conf-worker",
		}

		By("checking all machine configs for ones deployed by ztp")
		checkedConfigs := false
		machineConfigs, err := mco.ListMC(raninittools.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to list machine configs")

		for _, machineConfig := range machineConfigs {
			for _, machineConfigToCheck := range machineConfigsToCheck {
				if !strings.Contains(machineConfig.Object.Name, machineConfigToCheck) {
					continue
				}

				checkedConfigs = true

				glog.V(tsparams.LogLevel).Infof("Checking mc '%s' for annotation '%s'", machineConfig.Object.Name)

				annotation, ok := machineConfig.Object.Annotations[tsparams.ZtpGeneratedAnnotation]
				Expect(ok).To(BeTrue(), "Failed to find ZTP generated annotation")
				Expect(annotation).To(Equal("{}"), "ZTP generated annotation had the wrong val")
			}
		}

		By("checking if there were any matches")
		if !checkedConfigs {
			Skip("No matching machine configs were found")
		}
	})
})
