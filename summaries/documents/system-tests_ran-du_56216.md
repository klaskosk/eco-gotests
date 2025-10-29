# Test Case Summary for 56216

Test case 56216 is located in tests/system-tests/ran-du/tests/kernel-crash-kdump.go and is named "Trigger kernel crash to generate kdump vmcore".

## Goal

The goal of this test case is to verify that a kernel crash can be triggered on cluster nodes and that a `kdump vmcore` file is successfully generated as a result. It also ensures that the OpenShift API server remains available after the crash and reboot.

## Test Setup

Prior to the test case, the following changes are needed:

- Pull the `openshift-apiserver` deployment object to monitor its availability.

It does not require a git config set up.

## Test Steps

1. Retrieve the list of all nodes in the cluster.
2. For each node:
    a. Trigger a kernel crash on the node.
    b. Wait for the `openshift-apiserver` deployment to become available.
    c. Wait for the configured `RebootRecoveryTime` for cluster resources to reconcile their state.
    d. Assert that a `vmcore` dump file was generated in `/var/crash` on the node.
