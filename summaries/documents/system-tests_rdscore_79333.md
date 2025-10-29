# Test Case Summary for 79333

Test case 79333 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies pod-level bonded workloads on the same node and different PFs post graceful reboot".

## Goal

The goal of this test case is to verify that pod-level bonded workloads, deployed on the same node but utilizing different Physical Functions (PFs), remain functional and connected after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that pod-level bonding is configured on the same node with workloads deployed using different PFs.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyPodLevelBondWorkloadsOnSameNodeDifferentPFs`) to validate the functionality and connectivity of the bonded workloads.
2. This typically involves deploying a pod with a bonded interface on a specific node, configured to use different PFs.
3. The test then verifies network connectivity and performance for this bonded workload after a graceful reboot, confirming its resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable operation of pod-level bonded workloads on the same node with different PFs after a graceful reboot.
