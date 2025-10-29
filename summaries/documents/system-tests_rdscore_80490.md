# Test Case Summary for 80490

Test case 80490 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies pod-level bonded workloads after pod crashing".

## Goal

The goal of this test case is to verify the resilience and recovery of pod-level bonded workloads after a pod crash event.

## Test Setup

Prior to the test case, this test assumes a cluster with pod-level bonding configured. The test will simulate a pod crash scenario.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyPodLevelBondWorkloadsAfterPodCrashing` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that bonded workloads can recover and maintain network connectivity after a pod crashing, demonstrating resilience.
