# Test Case Summary for 82302

Test case 82302 is located in tests/cnf/ran/ptp/tests/ptp-events-and-metrics.go and is named "verifies phc2sys and ptp4l processes are UP".

## Goal

The goal of this test case is to validate that `phc2sys` and `ptp4l` processes remain `UP` after a PTP configuration change (specifically, adjusting and resetting the holdover timeout).

## Test Setup

Prior to the test case, the following conditions are asserted:

- A Prometheus API client is created.
- The PTP clocks are in a `LOCKED` state.

It does not require a git config set up.

## Test Steps

1.  Retrieve a map of node names to `NodeInfo` structs for nodes with PTP profiles using `profiles.GetNodeInfoMap`.
2.  For each node:
    a.  Get the first profile for the node using `nodeInfo.GetProfileByConfigPath` with "ptp4l.0.config".
    b.  Update the `HoldOverTimeout` for this profile to 60 seconds using `profiles.SetHoldOverTimeouts`, storing the original values for cleanup.
    c.  Defer a cleanup function to reset the `HoldOverTimeout` to original values using `profiles.ResetHoldOverTimeouts` and wait for the old timeouts to be applied using `profiles.WaitForOldHoldOverTimeouts`.
    d.  Wait for the new (60 seconds) `HoldOverTimeout` to be reflected in metrics using `profiles.WaitForHoldOverTimeouts`.
    e.  Reset the `HoldOverTimeout` for the profile to original values using `profiles.ResetHoldOverTimeouts`.
    f.  Wait for the old `HoldOverTimeout` values to be reflected in metrics using `profiles.WaitForOldHoldOverTimeouts`.
    g.  Ensure the process status is `UP` for both `phc2sys` and `ptp4l` using `metrics.AssertQuery` with a `ProcessStatusQuery` that includes both process types, the node name, and the config "ptp4l.0.config", with a timeout of 5 minutes.
3.  If no nodes with at least one profile were found, the test is skipped.
