# Test Case Summary for 54240

Test case 54240 (valid template scenario) is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-hub-templating.go and is named "should create the policy successfully with a valid template".

## Goal

The goal of this test case is to verify that a policy is created successfully when a valid ACM hub-side templating is used with TALM.

## Test Setup

Prior to the test case, the policies app's original source path is saved and then reset after the test. The test also ensures that the test namespace exists. It requires ZTP version 4.12 or later. If ZTP version is 4.16 or later, a different test path (`tsparams.ZtpTestPathTemplatingValid416`) is used. A secret named `tsparams.HubTemplatingSecretName` with opaque type and data `vlanQoS: MAo=` is created before the test and cleaned up afterwards.

It does require a git config set up such that the policies app can be updated via GitOps.

## Test Steps

1. Create a secret required for the valid template.
2. Call `setupHubTemplateTest` with the policies app and the appropriate valid ZTP test path (either `tsparams.ZtpTestPathTemplatingValid` or `tsparams.ZtpTestPathTemplatingValid416` based on ZTP version).
3. Pull the policy from the hub cluster.
4. Wait for the policy to reach a `Compliant` status.
