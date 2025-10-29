# Test Case Summary for 71848

Test case 71848 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies all NodeNetworkConfigurationPolicies are Available after ungraceful reboot".

## Goal

The goal of this test case is to verify that all NodeNetworkConfigurationPolicies (NNCPs) are in an "Available" state after an ungraceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and the cluster is in a state of recovery where NNCPs should be re-establishing their availability.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`nncp.WaitForAllNNCPsAvailable`) to wait for all NNCPs to become available.
2. It then verifies that no errors occurred during this wait and that all NNCPs have indeed reached the "Available" state.
3. The intent is to confirm the successful recovery and enforcement of network configuration policies after a disruptive event.
