# Test Case Summary for 82480

Test case 82480 is located in tests/cnf/ran/ptp/tests/ptp-events-and-metrics.go and is named "verifies all clocks are LOCKED".

## Goal

The goal of this test case is to verify that all PTP clocks on all nodes are in a `LOCKED` state.

## Test Setup

Prior to the test case, the following conditions are asserted:

- A Prometheus API client is created.
- The PTP clocks are in a `LOCKED` state.

It does not require a git config set up.

## Test Steps

1.  Ensure all clocks on all nodes are in a `LOCKED` state using `metrics.AssertQuery` with a stable duration of 10 seconds and a timeout of 5 minutes.
