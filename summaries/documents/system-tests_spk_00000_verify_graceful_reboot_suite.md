# Test Case Summary for Verify Graceful Reboot Suite

Test case Verify Graceful Reboot Suite is located in tests/system-tests/spk/internal/spkcommon/reboot.go and is named "Graceful reboot validation".

## Goal

The goal of this test suite is to validate the cluster's behavior after a graceful (soft) reboot, ensuring that all critical components like nodes, ClusterOperators, and deployments recover and become available.

## Test Setup

This test suite is called within the `Soft Reboot` describe block in `spk-suite.go`. Prior to this suite, the following setup is performed:

- The SPK backend UDP workload is set up using `spkcommon.SetupSPKBackendUDPWorkload()`.
- The SPK backend TCP workload is set up using `spkcommon.SetupSPKBackendWorkload()`.

After each test in the `Soft Reboot` suite, `spkcommon.ResetTMMReplicas()` is called to reset TMM replicas to 1.

It does not require a git config set up.

## Test Steps

1.  **Verifies graceful cluster reboot (ID 30021):** Calls `spkcommon.VerifySoftReboot()`. This function:
    - Skips the test if BMC details are not specified.
    - Gets a list of all nodes.
    - For each node, it cordons and drains the node.
    - Creates a BMC client for the node.
    - Performs a graceful reboot (power cycle) using the BMC client.
    - Waits for the node to become ready and then uncordons it.
2.  **Verifies all ClusterOperators are Available after graceful reboot (ID 72040):**
    - Waits for 3 minutes.
    - Checks if all ClusterOperators are available using `clusteroperator.WaitForAllClusteroperatorsAvailable()` within a 15-minute timeout.
3.  **Verifies all deployments are available after graceful reboot (ID 72041):** Calls `spkcommon.WaitAllDeploymentsAreAvailable()`. This function:
    - Lists all deployments in all namespaces.
    - Asserts that all deployments are in an "Available" state within a 25-minute timeout, with a 15-second polling interval.
