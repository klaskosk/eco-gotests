# Test Case Summary for 42736

Test case 42736 is located in tests/system-tests/ran-du/tests/hard-reboot.go and is named "Hard reboot nodes".

## Goal

The goal of this test case is to verify that a hard reboot of worker nodes does not disrupt the workload and that the cluster resources reconcile their state successfully. It also verifies SRIOV and PTP status after reboot.

## Test Setup

Prior to the test case, the following changes are needed:

- Prepare the workload by deleting any existing workload and then launching a new one using the shell method if specified.
- Wait for all deployments and statefulsets to become ready in the test workload namespace.

It does not require a git config set up.

## Test Steps

1. Retrieve the list of worker nodes.
2. For each configured hard reboot iteration:
    a. Hard reboot each worker node in the list.
    b. Wait for the configured `RebootRecoveryTime` for cluster resources to reconcile.
    c. Remove any pods in the `UnexpectedAdmissionError` state in the test workload namespace.
    d. Wait for all deployment and statefulset replicas to become ready in the test workload namespace.
    e. Retrieve the pod list and SRIOV networks with the `vfio-pci` driver.
    f. For each pod, assert that the number of devices under `/dev/vfio` is equal to or more than the pod's `vfio-pci` network attachments.
    g. If PTP is enabled, wait for 3 minutes and then validate the PTP status.
3. After all iterations, clean up the test workload resources by deleting the workload.
