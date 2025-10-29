# Test Case Summary for 78109

Test case 78109 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify eIP address from the list of defined does not used for the assigned pods in single eIP namespace with the wrong label".

## Goal

The goal of this test case is to verify that an EgressIP (eIP) address from a defined list is *not* used for assigned pods in a single eIP namespace when that namespace has an incorrect label.

## Test Setup

Prior to the test case, this test assumes that EgressIP is configured in the cluster with an eIP address list, and there are pods in a single eIP namespace with an incorrect namespace label that should prevent eIP assignment.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyEgressIPWrongNsLabel` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that eIP addresses are *not* assigned to pods in namespaces that do not meet the label selection criteria, thus validating the label-based assignment mechanism of EgressIP.
