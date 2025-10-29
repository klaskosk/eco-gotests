# Test Case Summary for 42738

Test case 42738 is located in tests/system-tests/ran-du/tests/soft-reboot.go and is named "Soft reboot nodes".

## Goal

The goal of this test case is to verify that soft rebooting cluster nodes does not disrupt the running workload. It ensures that after a soft reboot, cluster resources reconcile, the OpenShift API server remains available, and critical functionalities like SRIOV and PTP (if enabled) are restored and functional.

## Test Setup

Prior to the test case, the following changes are needed:

- Prepare the workload by deleting any existing workload and then launching a new one using the shell method if specified.
- Wait for all deployments and statefulsets to become ready in the test workload namespace.
- Pull the `openshift-apiserver` deployment object to monitor its availability.

It does not require a git config set up.

## Test Steps

1. Retrieve the list of all nodes in the cluster.
2. For each configured `SoftRebootIterations`:
    a. For each node in the cluster:
        i.   Perform a soft reboot on the node.
        ii.  Wait for the node to become unreachable.
        iii. Wait for the `openshift-apiserver` deployment to become "Available".
        iv.  Wait for the configured `RebootRecoveryTime` for cluster resources to reconcile their state.
        v.   Remove any pods in the `UnexpectedAdmissionError` state in the test workload namespace.
        vi.  Wait for all deployment replicas to become ready.
        vii. Wait for all statefulset replicas to become ready.
        viii. Retrieve the pod list and SRIOV networks with the `vfio-pci` driver.
        ix.  For each pod, assert that the number of devices under `/dev/vfio` is equal to or more than the pod's `vfio-pci` network attachments.
        x.   If PTP is enabled, wait for 3 minutes and then validate the PTP status, asserting it is in sync.
3. After all iterations, clean up the test workload resources by deleting the workload.
