package tests

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
)

var _ = Describe("PTP Node Reboot", Label(tsparams.Labels...), func() {
	// 59858 - verify the system returns to stability after reboot node
	It("should return to same stable status after ptp node soft reboot", reportxml.ID("59858"), func() {
	})

	// 59995 - Validates PTP consumer events after ptp node reboot
	It("validates PTP consumer events after ptp node reboot", reportxml.ID("59995"), func() {
	})
})
