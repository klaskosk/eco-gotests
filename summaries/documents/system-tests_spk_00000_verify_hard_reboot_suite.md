# Test Case Summary for Verify Hard Reboot Suite

Test case Verify Hard Reboot Suite is located in tests/system-tests/spk/internal/spkcommon/reboot.go and is named "Ungraceful reboot validation".

## Goal

The goal of this test suite is to validate the cluster's behavior after an ungraceful (hard) reboot, ensuring that all critical components like nodes, ClusterOperators, and deployments recover and become available.

## Test Setup

This test suite is called within the `Hard reboot` describe block in `spk-suite.go`. Prior to this suite, the following setup is performed:

- The SPK backend UDP workload is set up using `spkcommon.SetupSPKBackendUDPWorkload()`.
- The SPK backend TCP workload is set up using `spkcommon.SetupSPKBackendWorkload()`.

After each test in the `Hard reboot` suite, `spkcommon.ResetTMMReplicas()` is called to reset TMM replicas to 1.

It does not require a git config set up.

## Test Steps

1.  **Verifies ungraceful cluster reboot (ID 30020):** Calls `spkcommon.VerifyUngracefulReboot()`. This function:
    - Skips the test if BMC details are not specified.
    - Creates BMC clients for each node using the `SPKConfig.NodesCredentialsMap`.
    - For each node, it queries the power state (asserting "On") and then performs a `SystemForceReset` (hard reboot) using the BMC client.
    - Waits for all reboots to finish.
    - Calls `spkcommon.WaitAllNodesAreReady()` to wait for all nodes in the cluster to report a "Ready" state.
2.  **Verifies all ClusterOperators are Available after ungraceful reboot (ID 71868):**
    - Waits for 3 minutes.
    - Checks if all ClusterOperators are available using `clusteroperator.WaitForAllClusteroperatorsAvailable()` within a 15-minute timeout.
3.  **Removes all pods with UnexpectedAdmissionError:**
    - Waits for 3 minutes.
    - Lists pods in all namespaces with `status.phase=Failed`.
    - Iterates through the found pods and deletes those with `status.reason == "UnexpectedAdmissionError"`.
4.  **Verifies all deployments are available (ID 71872):** Calls `spkcommon.WaitAllDeploymentsAreAvailable()`. This function:
    - Lists all deployments in all namespaces.
    - Asserts that all deployments are in an "Available" state within a 25-minute timeout, with a 15-second polling interval.
