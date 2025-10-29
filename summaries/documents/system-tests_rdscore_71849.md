# Test Case Summary for 71849

Test case 71849 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies all NodeNetworkConfigurationPolicies are Available after soft reboot".

## Goal

The goal of this test case is to verify that all NodeNetworkConfigurationPolicies (NNCPs) are in an "Available" state after a graceful (soft) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and the cluster is in a state of recovery where NNCPs should be re-establishing their availability.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyAllNNCPsAreOK`) to check the availability of all NNCPs.
2. This typically involves monitoring the status of all `NodeNetworkConfigurationPolicy` resources in the cluster, waiting for them to report an "Available" status within a defined timeout.
3. The intent is to confirm the proper recovery and application of network configurations after a controlled cluster restart.
