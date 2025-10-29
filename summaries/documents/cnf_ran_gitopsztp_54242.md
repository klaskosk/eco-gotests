# Test Case Summary for 54242

Test case 54242 is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-policies-app.go and is named "Invalid time duration string for user override of policy intervals - should specify an invalid interval format and verify the app error".

## Goal

The goal of this test case is to verify that the Argo CD application reports an error when an invalid time duration string is used for overriding PGT policy intervals.

## Test Setup

Prior to the test case, the original policies app source is saved and then reset after the test. It requires ZTP version 4.10 or later.

It does require a git config set up such that the policies app can be updated via GitOps.

## Test Steps

1. Check if the ZTP test path for invalid intervals exists.
2. Update the Argo CD policies app with the invalid interval Git path (`tsparams.ZtpTestPathInvalidInterval`) and wait for synchronization.
3. Pull the Argo CD policies app.
4. Check the Argo CD conditions for the expected error message: "evaluationInterval.compliant 'time: invalid duration".
