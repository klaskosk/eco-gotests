# Test Case Summary for 66848

Test case 66848 is located in tests/cnf/ran/ptp/tests/ptp-events-and-metrics.go and is named "verifies phc2sys and ptp4l processes are UP".

## Goal

The goal of this test case is to verify that all `phc2sys` and `ptp4l` processes are in an `UP` state.

## Test Setup

Prior to the test case, the following conditions are asserted:

- A Prometheus API client is created.
- The PTP clocks are in a `LOCKED` state.

It does not require a git config set up.

## Test Steps

1.  Ensure all `phc2sys` and `ptp4l` processes are in an `UP` state using `metrics.AssertQuery` with a `ProcessStatusQuery` that includes both process types, and a timeout of 5 minutes.
