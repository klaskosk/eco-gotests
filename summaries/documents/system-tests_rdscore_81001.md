# Test Case Summary for 81001

Test case 81001 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies SR-IOV workloads on the same node and same SR-IOV network".

## Goal

The goal of this test case is to verify SR-IOV workloads deployed on the same node and utilizing the same SR-IOV network.

## Test Setup

Prior to the test case, no specific changes are mentioned as required. This test likely assumes a cluster with SR-IOV networks configured on the same node.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifySRIOVWorkloadsOnSameNode` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that SR-IOV enabled workloads can communicate when deployed on the same node and connected to the same SR-IOV network.
