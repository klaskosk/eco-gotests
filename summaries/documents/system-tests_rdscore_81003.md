# Test Case Summary for 81003

Test case 81003 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies SR-IOV workloads on different nodes and different SR-IOV networks".

## Goal

The goal of this test case is to verify SR-IOV workloads deployed on different nodes and utilizing different SR-IOV networks.

## Test Setup

Prior to the test case, no specific changes are mentioned as required. This test likely assumes a cluster with SR-IOV networks configured across multiple nodes.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifySRIOVWorkloadsOnDifferentNodesDifferentNet` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that SR-IOV enabled workloads can communicate when deployed on different nodes and connected to different SR-IOV networks.
