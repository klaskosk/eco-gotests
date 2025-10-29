# Test Case Summary for 79337

Test case 79337 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies pod-level bonded workloads on the different nodes and different PFs post graceful reboot".

## Goal

The goal of this test case is to verify that pod-level bonded workloads, deployed on different nodes and utilizing different Physical Functions (PFs), remain functional and connected after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that pod-level bonding is configured across different nodes with workloads deployed using different PFs.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyPodLevelBondWorkloadsOnDifferentNodesDifferentPFs`) to validate the functionality and connectivity of the bonded workloads.
2. This typically involves deploying pods with bonded interfaces across different nodes, each configured to use different PFs.
3. The test then verifies network connectivity and performance between these distributed bonded workloads after a graceful reboot, confirming their resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable operation of pod-level bonded workloads across different nodes with different PFs after a graceful reboot.
