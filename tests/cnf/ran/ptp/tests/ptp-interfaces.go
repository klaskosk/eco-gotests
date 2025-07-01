package tests

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/querier"
	. "github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ptp/internal/metrics"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

var _ = Describe("PTP Interfaces", Label(tsparams.LabelInterfaces), func() {
	var prometheusAPI prometheusv1.API

	BeforeEach(func() {
		By("creating a Prometheus API client")
		var err error
		prometheusAPI, err = querier.CreatePrometheusAPIForCluster(RANConfig.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to create Prometheus API client")

		By("ensuring clocks are locked before testing")
		err = metrics.AssertQuery(context.TODO(), prometheusAPI, metrics.ClockStateQuery{}, metrics.ClockStateLocked,
			metrics.AssertWithStableDuration[metrics.PtpClockState](10*time.Second),
			metrics.AssertWithTimeout[metrics.PtpClockState](5*time.Minute))
		Expect(err).ToNot(HaveOccurred(), "Failed to assert clock state is locked")
	})

	// 49742 - Validating events when slave interface goes down and up
	It("should generate events when slave interface goes down and up", reportxml.ID("49742"), func() {
		Fail("blocked by the ptp dynamic interface map")
	})
})
