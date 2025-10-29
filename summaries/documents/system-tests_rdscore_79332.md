# Test Case Summary for 79332

Test case 79332 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies pod-level bonded workloads on the same node and different PFs post hard reboot".

## Goal

The goal of this test case is to verify that pod-level bonded workloads, deployed on the same node but utilizing different Physical Functions (PFs), remain functional and connected after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that pod-level bonding is configured on the same node with workloads deployed using different PFs.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyPodLevelBondWorkloadsOnSameNodeDifferentPFs`) to validate the functionality and connectivity of the bonded workloads.
2. This typically involves performing network connectivity tests (e.g., ping, iperf) between bonded pods on the same node using different PFs, demonstrating their resilience post-reboot.
3. The intent is to confirm the proper functioning of pod-level bonded workloads after a disruptive event, especially in a single-node, multi-PF environment.
