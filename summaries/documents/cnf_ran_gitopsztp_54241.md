# Test Case Summary for 54241

Test case 54241 is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-policies-app.go and is named "User override of policy intervals - should specify new intervals and verify they were applied".

## Goal

The goal of this test case is to verify that a user can override the compliance and non-compliance intervals of a PolicyGeneratorTemplate (PGT) policy via GitOps and that the new intervals are correctly applied.

## Test Setup

Prior to the test case, the original policies app source is saved and then reset after the test. It requires ZTP version 4.10 or later.

It does require a git config set up such that the policies app can be updated via GitOps.

## Test Steps

1. Check if the ZTP test path for custom intervals exists.
2. Update the Argo CD policies app with the custom interval Git path (`tsparams.ZtpTestPathCustomInterval`) and wait for synchronization.
3. Wait for both the default policy (`tsparams.CustomIntervalDefaultPolicyName`) and the override policy (`tsparams.CustomIntervalOverridePolicyName`) to be created.
4. Validate the compliance and non-compliance intervals on the default policy, expecting them to be "1m".
5. Validate the compliance and non-compliance intervals on the overridden policy, expecting them to be "2m".
