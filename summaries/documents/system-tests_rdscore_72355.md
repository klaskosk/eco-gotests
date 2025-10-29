# Test Case Summary for 72355

Test case 72355 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies all policies are compliant after hard reboot".

## Goal

The goal of this test case is to verify that all cluster policies are compliant after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and the cluster is in a state of recovery where policies should be re-evaluated for compliance.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.ValidateAllPoliciesCompliant`) to check the compliance of all policies in the cluster.
2. It verifies that no errors occurred during the policy validation and that all policies are reported as compliant.
3. The overall intent is to confirm that the security and operational policies of the cluster are correctly enforced even after a disruptive event.
