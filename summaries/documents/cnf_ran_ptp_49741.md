# Test Case Summary for 49741

Test case 49741 is located in tests/cnf/ran/ptp/tests/ptp-events-and-metrics.go and is named "verifies FREERUN event received after adding a PHC offset".

## Goal

The goal of this test case is to verify that a `FREERUN` event is received after deliberately adjusting the PTP Hardware Clock (PHC) offset, and that the clock eventually returns to a `LOCKED` state.

## Test Setup

Prior to the test case, the following conditions are asserted:

- A Prometheus API client is created.
- The PTP clocks are in a `LOCKED` state.

It does not require a git config set up.

## Test Steps

1.  Retrieve a map of node names to `NodeInfo` structs for nodes with PTP profiles using `profiles.GetNodeInfoMap`.
2.  For each node:
    a.  Check for client interfaces on the node using `nodeInfo.GetInterfacesByClockType`; skip the node if none are found.
    b.  Group the client interfaces by NIC name using `iface.GroupInterfacesByNIC`.
    c.  Get the `cloud-event-consumer` pod for the node using `consumer.GetConsumerPodforNode`.
    d.  For each interface group (NIC):
        i.   Record the current time (`startTime`).
        ii.  Adjust the PHC by 5 milliseconds for the first interface in the group using `iface.AdjustPTPHardwareClock`.
        iii. Wait for a `PtpStateChange` event with `FREERUN` sync state using `events.WaitForEvent`.
        iv.  Record the current time (`startTime`).
        v.   Reset the PTP hardware clock for the interface using `iface.ResetPTPHardwareClock`.
        vi.  Wait for a `PtpStateChange` event with `LOCKED` sync state using `events.WaitForEvent`.
3.  If no nodes with at least one client interface were found, the test is skipped.
