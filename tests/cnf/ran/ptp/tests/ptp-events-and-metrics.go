package tests

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
)

var _ = Describe("PTP Events and Metrics", Ordered, ContinueOnFailure, Label(tsparams.Labels...), func() {
	BeforeEach(func() {
		By("ensuring clocks are locked before testing")
	})

	// we ensure the clocks are locked no less than 4 times before this test case runs, should probably remove
	It("should have [LOCKED] clock state in PTP metrics", reportxml.ID("82480"), func() {
	})

	// 66848 - Validate stability of the system before and after system manipulations
	It("should have the 'phc2sys' and 'ptp4l' processes in 'UP' state in PTP metrics", reportxml.ID("66848"), func() {
	})

	// 49741 - Change Offset Thresholds to min, max
	It("should change the slave clock state to free run after modify the offset threshold", reportxml.ID("49741"), func() {
	})

	// 82302 - Validating 'phc2sys' and 'ptp4l' processes state is 'UP' after ptp config change
	It("should have the 'phc2sys' and 'ptp4l' processes 'UP' after ptp config change", reportxml.ID("82302"), func() {
	})

})
