# Test Case Summary for 80967

Test case 80967 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies pod-level bonded workloads on the same node and same PF post hard reboot".

## Goal

The goal of this test case is to verify that pod-level bonded workloads, deployed on the same node and using the same Physical Function (PF), remain functional and connected after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that pod-level bonding is configured on the same node with workloads deployed using the same PF.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyPodLevelBondWorkloadsOnSameNodeSamePF`) to validate the functionality and connectivity of the bonded workloads.
2. This typically involves performing network connectivity tests (e.g., ping, iperf) between bonded pods on the same node sharing the same PF, demonstrating their resilience post-reboot.
3. The intent is to confirm the proper functioning of pod-level bonded workloads after a disruptive event, especially in a single-node, shared PF environment.
