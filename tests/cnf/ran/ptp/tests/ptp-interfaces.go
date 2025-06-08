package tests

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
)

var _ = Describe("PTP Interfaces", Label(tsparams.Labels...), func() {
	BeforeEach(func() {
		By("ensuring clocks are locked before testing")
	})

	// 49742 - Validating events when slave interface goes down and up
	It("should generate events when slave interface goes down and up", reportxml.ID("49742"), func() {
	})
})
