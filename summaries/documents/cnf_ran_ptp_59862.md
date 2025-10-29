# Test Case Summary for 59862

Test case 59862 is located in tests/cnf/ran/ptp/tests/ptp-process-restart.go and is named "should recover the phc2sys process after killing it".

## Goal

The goal of this test case is to validate that the `phc2sys` process recovers automatically after being killed, and that appropriate PTP events (`FREERUN` and `LOCKED`) are generated for `CLOCK_REALTIME`.

## Test Setup

Prior to the test case, the following conditions are asserted:

- A Prometheus API client is created.
- The PTP clocks are in a `LOCKED` state.

It does not require a git config set up.

## Test Steps

1.  Retrieve a map of node names to `NodeInfo` structs for nodes with PTP profiles using `profiles.GetNodeInfoMap`.
2.  For each node:
    a.  Get the `cloud-event-consumer` pod for the node using `consumer.GetConsumerPodforNode`.
    b.  Get the current PID of the `phc2sys` process using `processes.GetPID`.
    c.  Record the current time (`startTime`).
    d.  Kill the `phc2sys` process twice using `processes.KillPtpProcessMultipleTimes`.
    e.  Wait for a `OsClockSyncStateChange` event with `FREERUN` sync state on `CLOCK_REALTIME` using `events.WaitForEvent`.
    f.  Wait for a `OsClockSyncStateChange` event with `LOCKED` sync state on `CLOCK_REALTIME` using `events.WaitForEvent`.
    g.  Get the new PID of the `phc2sys` process and assert that it is different from the old PID, confirming a restart.
3.  If no nodes were found to run the test on, the test is skipped.
