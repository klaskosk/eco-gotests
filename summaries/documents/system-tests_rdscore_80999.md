# Test Case Summary for 80999

Test case 80999 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies SR-IOV workloads on different nodes and same SR-IOV network".

## Goal

The goal of this test case is to verify SR-IOV workloads deployed on different nodes but utilizing the same SR-IOV network.

## Test Setup

Prior to the test case, no specific changes are mentioned as required. This test likely assumes a cluster with SR-IOV networks configured across multiple nodes.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifySRIOVWorkloadsOnDifferentNodes` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that SR-IOV enabled workloads can communicate when deployed on different nodes but connected to the same SR-IOV network.
