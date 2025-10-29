# Test Case Summary for 74462

Test case 74462 is located in tests/system-tests/ran-du/tests/workload-guaranteed-force-delete.go and is named "Assert all pods recover after force deletion".

## Goal

The goal of this test case is to verify that all pods with a "Guaranteed" QoS class recover successfully after being forcefully deleted. This ensures the resilience of critical workloads in the event of unexpected pod terminations.

## Test Setup

Prior to the test case, the following changes are needed:

- Prepare the workload by deleting any existing workload and then launching a new one using the shell method if specified.
- Wait for all deployments and statefulsets to become ready in the test workload namespace.

It does not require a git config set up.

## Test Steps

1. For three iterations:
    a. List all pods in the test workload namespace.
    b. Identify and forcefully delete any pods with a "Guaranteed" QoS class.
    c. Assert that all pods in the test workload namespace recover and become ready within 2 minutes.
2. After all iterations, clean up the test workload resources by deleting the workload.
