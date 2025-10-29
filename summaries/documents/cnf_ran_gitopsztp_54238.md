# Test Case Summary for 54238

Test case 54238 is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-clusters-app.go and is named "User modification of klustletaddonconfig via gitops - should override the KlusterletAddonConfiguration and verify the change".

## Goal

The goal of this test case is to verify that a user can modify the KlusterletAddonConfiguration via GitOps and that the changes are applied and verified.

## Test Setup

Prior to the test case, the clusters app's original source path is saved and then reset after the test.

It does require a git config set up such that the clusters app can be updated via GitOps.

## Test Steps

1. Check if the ZTP test path for clusters app exists.
2. Update the clusters app Git path to `tsparams.ZtpTestPathClustersApp` and wait for synchronization.
3. Pull the `KlusterletAddonConfiguration` for `RANConfig.Spoke1Name`.
4. Wait for the `KlusterletAddonConfiguration` to have `SearchCollector` enabled.
