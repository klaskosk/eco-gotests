# Test Case Summary for 72566

Test case 72566 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify MACVLAN workload on different nodes".

## Goal

The goal of this test case is to verify MACVLAN workloads deployed on different nodes.

## Test Setup

Prior to the test case, no specific changes are mentioned as required. This test likely assumes a cluster with MACVLAN networks configured across multiple nodes.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyMacVlanOnDifferentNodes` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that MACVLAN-enabled workloads can communicate when deployed on different nodes.
