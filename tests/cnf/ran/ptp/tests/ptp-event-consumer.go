package tests

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	eventptp "github.com/redhat-cne/sdk-go/pkg/event/ptp"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/clients"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/daemonset"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/nodes"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/ptp"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/reportxml"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/querier"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/version"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/consumer"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/events"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/iface"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/metrics"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/profiles"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/ptpdaemon"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var _ = Describe("PTP Event Consumer", Label(tsparams.LabelEventConsumer), func() {
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

	// 64775 - Validate System is restored after POD restart/deletion
	It("should recover to stable state after delete PTP daemon pod", reportxml.ID("64775"), func() {
		testRanAtLeastOnce := false

		nodeInfoMap, err := profiles.GetNodeInfoMap(RANConfig.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to get node info map")

		for _, nodeInfo := range nodeInfoMap {
			testRanAtLeastOnce = true

			By("getting the PTP daemon pod for node " + nodeInfo.Name)
			ptpDaemonPod, err := ptpdaemon.GetPtpDaemonPodOnNode(RANConfig.Spoke1APIClient, nodeInfo.Name)
			Expect(err).ToNot(HaveOccurred(), "Failed to get PTP daemon pod for node %s", nodeInfo.Name)

			By("deleting the PTP daemon pod and waiting until it is deleted")
			startTime := time.Now()
			_, err = ptpDaemonPod.DeleteAndWait(5 * time.Minute)
			Expect(err).ToNot(HaveOccurred(), "Failed to delete PTP daemon pod for node %s", nodeInfo.Name)

			By("waiting for the PTP daemonset to be ready again")
			ptpDaemonset, err := daemonset.Pull(
				RANConfig.Spoke1APIClient, ranparam.LinuxPtpDaemonsetName, ranparam.PtpOperatorNamespace)
			Expect(err).ToNot(HaveOccurred(), "Failed to pull PTP daemon set")

			ready := ptpDaemonset.IsReady(5 * time.Minute)
			Expect(ready).To(BeTrue(), "Failed to wait for PTP daemon set to be ready")

			By("waiting up to 10 minutes for metrics to be locked and stable for 1 minute")
			err = metrics.AssertQuery(context.TODO(), prometheusAPI, metrics.ClockStateQuery{}, metrics.ClockStateLocked,
				metrics.AssertWithStableDuration[metrics.PtpClockState](1*time.Minute),
				metrics.AssertWithTimeout[metrics.PtpClockState](10*time.Minute))
			Expect(err).ToNot(HaveOccurred(), "Failed to assert clock state is locked and stable after pod restart")

			By("waiting up to 10 minutes since startTime for the locked event on the node")
			eventPod, err := consumer.GetConsumerPodforNode(RANConfig.Spoke1APIClient, nodeInfo.Name)
			Expect(err).ToNot(HaveOccurred(), "Failed to get event pod for node %s", nodeInfo.Name)

			filter := events.All(
				events.IsType(eventptp.PtpStateChange),
				events.HasValue(events.WithSyncState(eventptp.LOCKED)),
			)
			err = events.WaitForEvent(eventPod, startTime, 10*time.Minute, filter)
			Expect(err).ToNot(HaveOccurred(), "Failed to wait for locked event on node %s", nodeInfo.Name)
		}

		if !testRanAtLeastOnce {
			Skip("No nodes found to run the test on")
		}
	})

	// 54245 - Test event publisher/subscription via HTTP
	It("validates HTTP PTP events via consumer", reportxml.ID("54245"), func() {
		Skip("the test for 49741 already covers this")
	})

	// 59996 - Functional_Case: Working without consumer
	It("validates the system is fully functional after removing consumer", reportxml.ID("59996"), func() {
		Skip("not a realistic customer use case")
	})

	// 82218 - Validates the consumer events after ptpoperatorconfig api version is modified
	It("validates the consumer events after ptpoperatorconfig api version is modified", reportxml.ID("82218"), func() {
		By("checking if the PTP version is within the 4.16-4.18 range")
		inRange, err := version.IsVersionStringInRange(RANConfig.Spoke1OperatorVersions[ranparam.PTP], "4.16", "4.18")
		Expect(err).ToNot(HaveOccurred(), "Failed to check PTP version range")

		if !inRange {
			Skip("PTP version is not within the 4.16-4.18 range, skipping test")
		}

		By("cleaning up all consumers using consumer.CleanupConsumersOnWorkers")
		err = consumer.CleanupConsumersOnWorkers(RANConfig.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to cleanup consumers on workers")

		By("retrieving the current API version from the PTP Operator Config")
		ptpOperatorConfig, err := ptp.PullPtpOperatorConfig(RANConfig.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull PTP operator config")

		originalAPIVersion := ptpOperatorConfig.Definition.Spec.EventConfig.ApiVersion

		By("modifying the ptpEventConfig/apiVersion field in the PTP Operator Config")
		var newAPIVersion string
		if originalAPIVersion == "2.0" {
			newAPIVersion = "1.0"
		} else {
			newAPIVersion = "2.0"
		}

		ptpOperatorConfig.Definition.Spec.EventConfig.ApiVersion = newAPIVersion
		ptpOperatorConfig, err = ptpOperatorConfig.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update PTP operator config with new API version")

		By("waiting for the changes to propagate")
		time.Sleep(1 * time.Minute)

		By("verifying that all PTP clocks are in a LOCKED state")
		err = metrics.AssertQuery(context.TODO(), prometheusAPI, metrics.ClockStateQuery{}, metrics.ClockStateLocked,
			metrics.AssertWithStableDuration[metrics.PtpClockState](1*time.Minute),
			metrics.AssertWithTimeout[metrics.PtpClockState](10*time.Minute))
		Expect(err).ToNot(HaveOccurred(), "Failed to assert all clocks are locked after API version change")

		By("redeploying all the consumers")
		err = consumer.DeployConsumersOnWorkers(RANConfig.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to deploy consumers on workers")

		By("verifying that we see a PtpStateChange to LOCKED containing iface.Master")
		verifyPtpLockedEventOnWorkers(RANConfig.Spoke1APIClient)

		By("cleaning up all consumers")
		err = consumer.CleanupConsumersOnWorkers(RANConfig.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to cleanup consumers on workers")

		By("restoring the original PTP Operator Config")
		ptpOperatorConfig.Definition.Spec.EventConfig.ApiVersion = originalAPIVersion
		_, err = ptpOperatorConfig.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to restore original PTP operator config")

		By("waiting for the PTP clocks to return to a LOCKED state")
		err = metrics.AssertQuery(context.TODO(), prometheusAPI, metrics.ClockStateQuery{}, metrics.ClockStateLocked,
			metrics.AssertWithStableDuration[metrics.PtpClockState](1*time.Minute),
			metrics.AssertWithTimeout[metrics.PtpClockState](10*time.Minute))
		Expect(err).ToNot(HaveOccurred(), "Failed to assert all clocks are locked after restoring original config")

		By("redeploying all the consumers again")
		err = consumer.DeployConsumersOnWorkers(RANConfig.Spoke1APIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to redeploy consumers on workers")

		By("verifying that we see a PtpStateChange to LOCKED containing iface.Master again")
		verifyPtpLockedEventOnWorkers(RANConfig.Spoke1APIClient)
	})
})

func verifyPtpLockedEventOnWorkers(client *clients.Settings) {
	workerNodes, err := nodes.List(client, metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set(RANConfig.WorkerLabelMap)).String(),
	})
	Expect(err).ToNot(HaveOccurred(), "Failed to list worker nodes")

	for _, worker := range workerNodes {
		workerNodeName := worker.Definition.Name

		By("getting the event pod for node " + workerNodeName)

		eventPod, err := consumer.GetConsumerPodforNode(client, workerNodeName)
		Expect(err).ToNot(HaveOccurred(), "Failed to get event pod for node %s", workerNodeName)

		By("waiting for PtpStateChange to LOCKED containing iface.Master on node " + workerNodeName)

		filter := events.All(
			events.IsType(eventptp.PtpStateChange),
			events.HasValue(events.WithSyncState(eventptp.LOCKED), events.ContainingResource(string(iface.Master))),
		)
		err = events.WaitForEvent(eventPod, time.Now(), 5*time.Minute, filter)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for locked event with iface.Master on node %s", workerNodeName)
	}
}
