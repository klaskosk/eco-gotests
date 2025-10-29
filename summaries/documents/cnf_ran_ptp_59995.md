# Test Case Summary for 59995

Test case 59995 is located in tests/cnf/ran/ptp/tests/ptp-node-reboot.go and is named "validates PTP consumer events after ptp node reboot".

## Goal

The goal of this test case is to validate that the PTP consumer generates the `LOCKED` event after a node reboot.

## Test Setup

This test case reuses the setup from test case 59858, which involves:

1.  Determining if the cluster is SNO.
2.  Selecting a PTP node to reboot.
3.  Recording the reboot time.
4.  Soft rebooting the selected node.
5.  Waiting for the node to recover (either SNO-specific recovery or general node/pod health checks).

It does not require a git config set up.

## Test Steps

1.  Get the `cloud-event-consumer` pod for the rebooted node using `consumer.GetConsumerPodforNode`.
2.  Wait for a `PtpStateChange` event with `LOCKED` sync state, containing the resource `iface.Master`, to be reported by the event pod, using `events.WaitForEvent` with the `rebootTime` as the start time and a 5-minute timeout.
