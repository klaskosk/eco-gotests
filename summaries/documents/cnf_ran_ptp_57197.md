# Test Case Summary for 57197

Test case 57197 is located in tests/cnf/ran/ptp/tests/ptp-process-restart.go and is named "ensures ptp4l is restarted after killing ptp4l unrelated to phc2sys".

## Goal

The goal of this test case is to verify that the `ptp4l` process restarts automatically after being killed (specifically, a `ptp4l` process unrelated to `phc2sys`), that the `phc2sys` process is not affected, and that appropriate PTP events (`FREERUN` and `LOCKED`) are generated.

## Test Setup

Prior to the test case, the following conditions are asserted:

- A Prometheus API client is created.
- The PTP clocks are in a `LOCKED` state.

It does not require a git config set up.

## Test Steps

1.  Retrieve a map of node names to `NodeInfo` structs for nodes with PTP profiles using `profiles.GetNodeInfoMap`.
2.  For each node:
    a.  Check if there are at least two PTP profiles on the node; skip the node if not.
    b.  Update the `HoldOverTimeout` for all profiles on the node to 180 seconds using `profiles.SetHoldOverTimeouts`, storing the original values for cleanup.
    c.  Defer a cleanup function to reset the `HoldOverTimeout` to original values using `profiles.ResetHoldOverTimeouts` and wait for the old timeouts to be applied using `profiles.WaitForOldHoldOverTimeouts`.
    d.  Wait for the new (180 seconds) `HoldOverTimeout` to be reflected in metrics using `profiles.WaitForHoldOverTimeouts`.
    e.  Get the `cloud-event-consumer` pod for the node using `consumer.GetConsumerPodforNode`.
    f.  Get the PIDs of `ptp4l` processes that are *not* related to `phc2sys` using `processes.GetPtp4lPIDsByRelatedProcess`.
    g.  Get the current PID of the `phc2sys` process using `processes.GetPID`.
    h.  Record the current time (`startTime`).
    i.  Kill the first identified `ptp4l` process (unrelated to `phc2sys`) using `processes.KillProcessByPID`.
    j.  Wait for a `PtpStateChange` event with `FREERUN` sync state on `iface.Master` to be received using `events.WaitForEvent`.
    k.  Wait for a `PtpStateChange` event with `LOCKED` sync state on `iface.Master` to be received using `events.WaitForEvent`.
    l.  Ensure the `phc2sys` process PID has not changed, confirming it was unaffected.
    m.  Ensure a new `ptp4l` process has started (new PID) and the killed PID is no longer present using `processes.GetPtp4lPIDsByRelatedProcess`.
    n.  Reset the `HoldOverTimeout` for all profiles to original values using `profiles.ResetHoldOverTimeouts`.
    o.  Wait for the old `HoldOverTimeout` values to be reflected in metrics using `profiles.WaitForOldHoldOverTimeouts`.
3.  If no nodes were found to run the test on, the test is skipped.
