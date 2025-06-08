package tests

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
)

var _ = Describe("PTP Process Restart", Label(tsparams.Labels...), func() {
	// 59862 - validate phc2sys process restart after killing that process
	It("should recover the phc2sys process after killing it", reportxml.ID("59862"), func() {
	})

	// 57197 - Ptp4l restart - single process - Dual Nic
	It("should create a new ptp4l process after killing a ptp4l process that is not related to the "+
		"phc2sy process", reportxml.ID("57197"), func() {
	})
})
