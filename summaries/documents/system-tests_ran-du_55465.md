# Test Case Summary for 55465

Test case 55465 is located in tests/system-tests/ran-du/tests/launch-workload.go and is named "Assert all pods are ready".

## Goal

The goal of this test case is to verify that a deployed workload's pods, deployments, and statefulsets all reach a ready state, and that PTP synchronization is maintained (if enabled) after the workload deployment.

## Test Setup

Prior to the test case, the following changes are needed:

- Prepare the workload by deleting any existing workload and then launching a new one using the shell method if specified.
- Wait for all deployments and statefulsets to become ready in the test workload namespace.

It does not require a git config set up.

## Test Steps

1. Assert that all pods in the test workload namespace become ready within the default timeout.
2. If PTP is enabled, wait for 3 minutes and then check PTP status, asserting that it is in sync.
3. After the test, clean up the test workload resources by deleting the workload.
