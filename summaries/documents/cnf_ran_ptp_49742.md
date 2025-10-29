# Test Case Summary for 49742

Test case 49742 is located in tests/cnf/ran/ptp/tests/ptp-interfaces.go and is named "should generate events when slave interface goes down and up".

## Goal

The goal of this test case is to verify that PTP events are generated when a slave interface goes down and up. Specifically, it asserts that `HOLDOVER` and `FREERUN` events are generated when an interface goes down, and a `LOCKED` event is generated when the interface comes back up.

## Test Setup

Prior to the test case, the following conditions are asserted:

- A Prometheus API client is created.
- The PTP clocks are in a `LOCKED` state.

It does not require a git config set up.

## Test Steps

1.  Retrieve a map of node names to `NodeInfo` structs for nodes with PTP profiles using `profiles.GetNodeInfoMap`.
2.  For each node, identify receiver interfaces (slave interfaces) using `nodeInfo.GetInterfacesByClockType`.
3.  Determine the egress interface name for the node using `iface.GetEgressInterfaceName` and skip testing if a receiver interface belongs to the egress NIC.
4.  Group the receiver interfaces by their NIC names using `iface.GroupInterfacesByNIC`.
5.  For each group of interfaces:
    a.  Get the `cloud-event-consumer` pod for the node using `consumer.GetConsumerPodforNode`.
    b.  Defer a cleanup function to ensure all interfaces in the group are set to `up` at the end of the test, even if it fails, using `iface.SetInterfaceStatus`.
    c.  Set all interfaces in the current group to `down` using `iface.SetInterfaceStatus`.
    d.  Wait for a `PtpStateChange` event with `HOLDOVER` sync state on the relevant NIC using `events.WaitForEvent`.
    e.  Wait for a `PtpStateChange` event with `FREERUN` sync state on the relevant NIC using `events.WaitForEvent`.
    f.  Assert that the `ClockStateQuery` for the interface group on that node shows a `FREERUN` metric using `metrics.AssertQuery`.
    g.  Set all interfaces in the group to `up` using `iface.SetInterfaceStatus`.
    h.  Wait for a `PtpStateChange` event with `LOCKED` sync state on the relevant NIC using `events.WaitForEvent`.
    i.  Assert that all `ClockStateQuery` metrics are `LOCKED` using `metrics.AssertQuery`.
6.  If no interfaces were found to test, the test is skipped.
