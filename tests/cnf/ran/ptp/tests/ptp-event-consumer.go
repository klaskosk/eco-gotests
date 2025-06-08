package tests

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
)

var _ = Describe("PTP Event Consumer", Label(tsparams.Labels...), func() {
	// 64775 - Validate System is restored after POD restart/deletion
	It("should recover to stable state after delete PTP daemon pod", reportxml.ID("64775"), func() {
	})

	// 54245 - Test event publisher/subscription via HTTP
	It("validates HTTP PTP events via consumer", reportxml.ID("54245"), func() {
	})

	// 59996 - Functional_Case: Working without consumer
	It("validates the system is fully functional after removing consumer", reportxml.ID("59996"), func() {
	})

	// 82218 - Validates the consumer events after ptpoperatorconfig api version is modified
	It("validates the consumer events after ptpoperatorconfig api version is modified", reportxml.ID("82218"), func() {
	})
})
