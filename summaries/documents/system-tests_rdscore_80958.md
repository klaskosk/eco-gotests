# Test Case Summary for 80958

Test case 80958 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies pod-level bonded workloads on the same node and same PF".

## Goal

The goal of this test case is to verify the functionality of pod-level bonded workloads when deployed on the same node and utilizing the same Physical Function (PF).

## Test Setup

Prior to the test case, this test assumes a cluster with pod-level bonding configured, and that bonded workloads are deployed on the same node using the same PF.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyPodLevelBondWorkloadsOnSameNodeSamePF` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that pod-level bonded workloads function correctly and provide network redundancy or increased bandwidth in this specific configuration.
