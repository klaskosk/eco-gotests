package tests

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/reportxml"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/querier"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/metrics"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
)

var _ = Describe("PTP OC Two Port", Label(tsparams.LabelOC2Port), func() {
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

	// 80963 - Verifies 2-Port OC HA Failover when active port goes down
	It("verifies 2-port oc ha failover when active port goes down", reportxml.ID("80963"), func() {
	})
	// 80964 - Verifies 2-Port OC HA Holdover & Freerun when both ports go down
	It("verifies 2-port oc ha holdover & freerun when both ports go down", reportxml.ID("80964"), func() {
	})
	// 82012 - Verifies 2-Port OC HA passive interface recovery
	It("verifies 2-port oc ha passive interface recovery", reportxml.ID("82012"), func() {
	})
})
