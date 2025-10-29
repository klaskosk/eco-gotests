# Test Case Summary for 81002

Test case 81002 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies SR-IOV workloads on same node and different SR-IOV networks".

## Goal

The goal of this test case is to verify SR-IOV workloads deployed on the same node but utilizing different SR-IOV networks.

## Test Setup

Prior to the test case, no specific changes are mentioned as required. This test likely assumes a cluster with SR-IOV networks configured and available for deployment.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifySRIOVWorkloadsOnSameNodeDifferentNet` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that SR-IOV enabled workloads can communicate when deployed on the same node but connected to different SR-IOV networks.
