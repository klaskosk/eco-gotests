package tests

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/ptp"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/reportxml"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/querier"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/gnss"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/metrics"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/profiles"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/ptpdaemon"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
)

var _ = Describe("PTP GNSS with NTP Fallback", Label(tsparams.LabelNTPFallback), func() {
	var prometheusAPI prometheusv1.API

	BeforeEach(func() {
		By("creating a Prometheus API client")
		var err error
		prometheusAPI, err = querier.CreatePrometheusAPIForCluster(RANConfig.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to create Prometheus API client")

		By("ensuring clocks are locked before testing")
		err = metrics.AssertQuery(context.TODO(), prometheusAPI, metrics.ClockStateQuery{}, metrics.ClockStateLocked,
			metrics.AssertWithStableDuration(10*time.Second),
			metrics.AssertWithTimeout(5*time.Minute))
		Expect(err).ToNot(HaveOccurred(), "Failed to assert clock state is locked")
	})

	AfterEach(func() {
		By("ensuring clocks are locked after testing")
		err := metrics.AssertQuery(context.TODO(), prometheusAPI, metrics.ClockStateQuery{}, metrics.ClockStateLocked,
			metrics.AssertWithStableDuration(10*time.Second),
			metrics.AssertWithTimeout(5*time.Minute))
		Expect(err).ToNot(HaveOccurred(), "Failed to assert clock state is locked")
	})

	When("updating the PtpConfig to adjust ts2phc holdover", func() {
		var (
			savedPtpConfigs []*ptp.PtpConfigBuilder
			nodeInfoMap     profiles.NodeInfoMap
		)

		BeforeEach(func() {
			By("saving PtpConfigs before test")
			var err error
			savedPtpConfigs, err = profiles.SavePtpConfigs(RANConfig.Spoke1APIClient)
			Expect(err).ToNot(HaveOccurred(), "Failed to save PtpConfigs")

			By("getting node info map")
			nodeInfoMap, err = profiles.GetNodeInfoMap(RANConfig.Spoke1APIClient)
			Expect(err).ToNot(HaveOccurred(), "Failed to get node info map")
		})

		AfterEach(func() {
			// If the test did not fail, the PtpConfigs were restored by the test itself.
			if !CurrentSpecReport().Failed() {
				return
			}

			By("restoring PtpConfigs")
			startTime := time.Now()
			changedProfiles, err := profiles.RestorePtpConfigs(RANConfig.Spoke1APIClient, savedPtpConfigs)
			Expect(err).ToNot(HaveOccurred(), "Failed to restore PtpConfigs")

			By("getting node names for changed profiles")
			nodeNames := nodeInfoMap.GetNodesWithProfiles(changedProfiles)

			By("waiting for profile load on nodes")
			err = ptpdaemon.WaitForProfileLoadOnNodes(RANConfig.Spoke1APIClient, nodeNames, startTime, 5*time.Minute)
			Expect(err).ToNot(HaveOccurred(), "Failed to wait for profile load on nodes")
		})

		// 85904 - Successful fallback to NTP when GNSS sync lost
		It("successfully falls back to NTP when GNSS sync lost", reportxml.ID("85904"), func() {
			testActuallyRan := false

			By("getting node info map")
			nodeInfoMap, err := profiles.GetNodeInfoMap(RANConfig.Spoke1APIClient)
			Expect(err).ToNot(HaveOccurred(), "Failed to get node info map")

			for nodeName, nodeInfo := range nodeInfoMap {
				if nodeInfo.Counts[profiles.ProfileTypeNTPFallback] == 0 {
					continue
				}

				testActuallyRan = true

				By("getting the u-blox protocol version")
				ntpFallbackProfiles := nodeInfo.GetProfilesByTypes(profiles.ProfileTypeNTPFallback)
				Expect(ntpFallbackProfiles).ToNot(BeEmpty(), "No NTP fallback profile found for node %s", nodeName)

				ntpFallbackProfile, err := ntpFallbackProfiles[0].PullProfile(RANConfig.Spoke1APIClient)
				Expect(err).ToNot(HaveOccurred(), "Failed to pull NTP fallback profile for node %s", nodeName)

				protocolVersion, err := gnss.GetUbloxProtocolVersion(ntpFallbackProfile)
				Expect(err).ToNot(HaveOccurred(), "Failed to get u-blox protocol version for node %s", nodeName)

				By("setting the ts2phc holdover to 10 seconds")
				oldProfile, err := profiles.UpdateTS2PHCHoldover(RANConfig.Spoke1APIClient, ntpFallbackProfiles[0], 10)
				Expect(err).ToNot(HaveOccurred(), "Failed to update ts2phc holdover for node %s", nodeName)

				By("simulating GNSS sync loss")
				err = gnss.SimulateSyncLoss(RANConfig.Spoke1APIClient, nodeName, protocolVersion)
				Expect(err).ToNot(HaveOccurred(), "Failed to simulate GNSS sync loss for node %s", nodeName)

				DeferCleanup(func() {
					if !CurrentSpecReport().Failed() {
						return
					}

					By("restoring GNSS sync")
					err = gnss.SimulateSyncRecovery(RANConfig.Spoke1APIClient, nodeName, protocolVersion)
					Expect(err).ToNot(HaveOccurred(), "Failed to simulate GNSS sync recovery for node %s", nodeName)
				})

				By("restoring the ts2phc holdover")
				err = profiles.RestoreProfileToConfig(RANConfig.Spoke1APIClient, ntpFallbackProfiles[0], oldProfile)
				Expect(err).ToNot(HaveOccurred(), "Failed to restore ts2phc holdover for node %s", nodeName)

				By("waiting for profiles to be loaded")
			}

			if !testActuallyRan {
				Skip("No receiver interfaces found for any node")
			}
		})

		// 85905 - Failed fallback to NTP when GNSS sync lost
		It("fails to fall back to NTP when NTP server unreachable", reportxml.ID("85905"), func() {})

		// 85906 - Ensure system clock is within 1.5 ms for entire holdover
		It("verifies system clock is within 1.5 ms for entire ts2phc holdover", reportxml.ID("85906"), func() {})
	})
})
