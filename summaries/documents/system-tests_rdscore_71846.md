# Test Case Summary for 71846

Test case 71846 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies all NodeNetworkConfigurationPolicies are Available".

## Goal

The goal of this test case is to verify that all NodeNetworkConfigurationPolicies (NNCPs) are in an "Available" state.

## Test Setup

Prior to the test case, this test assumes that NodeNetworkConfigurationPolicies have been created and applied in the cluster.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyAllNNCPsAreOK` to perform the verification. The detailed steps are within this helper function, but the overall intent is to check the status of all NNCPs and ensure they are successfully applied and available.
