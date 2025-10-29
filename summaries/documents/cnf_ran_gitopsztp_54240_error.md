# Test Case Summary for 54240

Test case 54240 (error scenario) is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-hub-templating.go and is named "should report an error for using autoindent function where not allowed".

## Goal

The goal of this test case is to verify that the system correctly reports an error when an unsupported `autoindent` function is used in ACM hub-side templating with TALM.

## Test Setup

Prior to the test case, the policies app's original source path is saved and then reset after the test. The test also ensures that the test namespace exists. It requires ZTP version 4.12 or later.

It does require a git config set up such that the policies app can be updated via GitOps.

## Test Steps

1. Call `setupHubTemplateTest` with the policies app and `tsparams.ZtpTestPathTemplatingAutoIndent`.
2. Validate that TALM reports a policy error by asserting the TALM pod log contains "policy has hub template error".
3. Pull the hub-side templating policy.
4. Validate the specific error using the policy message, expecting it to contain "wrong type for value; expected string; got int".
