# Test Case Summary for 74522

Test case 74522 is located in tests/system-tests/ran-du/tests/stability-no-workload.go and is named "StabilityNoWorkload".

## Goal

The goal of this test case is to verify the stability of the cluster over a period of time without any active workload. It collects and verifies metrics related to PTP status, policy changes, pod restarts in critical namespaces (openshift-etcd, openshift-apiserver), and tuned restarts.

## Test Setup

Prior to the test case, the following changes are needed:

- Clean up any existing test workload resources.
- Fetch the OpenShift Cluster name.

It does not require a git config set up.

## Test Steps

1. Initialize output directories and file paths for storing stability metrics.
2. For the configured `StabilityNoWorkloadDurMins` duration, at intervals of `StabilityNoWorkloadIntMins`:
    a. If PTP is enabled, save the PTP status.
    b. If `StabilityPoliciesCheck` is enabled, save the policy status.
    c. For namespaces "openshift-etcd" and "openshift-apiserver", save pod restarts.
    d. Save tuned restarts.
3. After the collection period, perform final checks:
    a. If `StabilityPoliciesCheck` is enabled, verify that there are no policy changes.
    b. For namespaces "openshift-etcd" and "openshift-apiserver", verify no unexpected pod restarts.
    c. If PTP is enabled, verify PTP status changes.
    d. Verify no tuned restarts.
4. Assert that no stability errors occurred during the test.
