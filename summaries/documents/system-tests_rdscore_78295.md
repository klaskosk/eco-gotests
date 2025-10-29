# Test Case Summary for 78295

Test case 78295 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies pod-level bonded workloads on the different nodes and different PFs".

## Goal

The goal of this test case is to verify the functionality of pod-level bonded workloads when deployed on different nodes and utilizing different Physical Functions (PFs).

## Test Setup

Prior to the test case, this test assumes a cluster with pod-level bonding configured, and that bonded workloads are deployed on different nodes using different PFs.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyPodLevelBondWorkloadsOnDifferentNodesDifferentPFs` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that pod-level bonded workloads function correctly and provide network redundancy or increased bandwidth in this specific configuration.
