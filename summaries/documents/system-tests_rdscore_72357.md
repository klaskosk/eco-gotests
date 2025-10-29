# Test Case Summary for 72357

Test case 72357 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies all policies are compliant after soft reboot".

## Goal

The goal of this test case is to verify that all cluster policies are compliant after a graceful (soft) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and the cluster is in a state of recovery where policies should be re-evaluated for compliance.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.ValidateAllPoliciesCompliant`) to check the compliance of all policies in the cluster.
2. This typically involves iterating through all defined policies and verifying their status against expected compliance criteria.
3. The intent is to confirm that the cluster's policy enforcement mechanisms are fully operational and correctly applied after a controlled cluster restart.
