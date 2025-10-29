# Test Case Summary for 60619

Test case 60619 is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-clusters-app.go and is named "Image creation fails when NMstateConfig CR is empty - should not have NMStateConfig CR when nodeNetwork section not in siteConfig".

## Goal

The goal of this test case is to verify that when the `nodeNetwork` section is not present in the `siteConfig`, there should be no `NMStateConfig` CR on the hub.

## Test Setup

Prior to the test case, the clusters app's original source path is saved and then reset after the test.

It does require a git config set up such that the clusters app can be updated via GitOps.

## Test Steps

1. Check if the ZTP test path for removing NMState exists.
2. Check if the `NMStateConfig` CR exists on the hub. It should not be empty.
3. Update the clusters app Git path to `tsparams.ZtpTestPathRemoveNmState` and wait for synchronization.
4. Validate that the `NMStateConfig` is no longer present on the hub.
