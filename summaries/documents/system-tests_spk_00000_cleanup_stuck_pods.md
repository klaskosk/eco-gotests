# Test Case Summary for system-tests_spk_tests_no_id_cleanup_stuck_pods

Test case with no explicit ID, but labeled as "Removes stuck SPK pods", is located in tests/system-tests/spk/tests/spk-suite.go.

## Goal

The goal of this test case is to remove any SPK pods that are in a "Pending" phase with a "ContainerCreating" state, typically after a hard reboot.

## Test Setup

This test case is part of the "Hard reboot" describe block and runs after a hard reboot has been performed. It assumes there might be pods stuck in a `ContainerCreating` state in the `SPKConfig.SPKDataNS` and `SPKConfig.SPKDnsNS` namespaces.

It does not require a git config set up.

## Test Steps

1. The `cleanupStuckContainerPods` function is called for `SPKConfig.SPKDataNS`.
2. Inside `cleanupStuckContainerPods`, it lists pods in the specified namespace with `status.phase=Pending`.
3. It then iterates through the found pods and deletes each one that is in a "ContainerCreating" state, waiting for its termination.
4. Steps 1-3 are repeated for `SPKConfig.SPKDnsNS`.
