# Test Case Summary for 72567

Test case 72567 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify MACVLAN workloads on the same node".

## Goal

The goal of this test case is to verify MACVLAN workloads deployed on the same node.

## Test Setup

Prior to the test case, no specific changes are mentioned as required. This test likely assumes a cluster with MACVLAN networks configured on the same node.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyMacVlanOnSameNode` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that MACVLAN-enabled workloads can communicate when deployed on the same node.
