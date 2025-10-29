# Test Case Summary for 75383

Test case 75383 is located in tests/cnf/ran/gitopsztp/tests/ztp-siteconfig-failover.go and is named "Verify siteconfig operator's recovery mechanism by referencing non-existent extra manifests configmap custom resource".

## Goal

The goal of this test case is to verify that the siteconfig operator's recovery mechanism works correctly when a `ClusterInstance` CR references a non-existent extra manifests configmap, and that a proper validation error is reported.

## Test Setup

Prior to the test case, the original `clustersApp` source is saved and then reset in `BeforeEach` and `AfterEach` blocks respectively. It requires ZTP version 4.17 or later.

It does require a git config set up such that the clusters app can be updated via GitOps.

## Test Steps

1. Check if the non-existent extra manifests configmap reference Git path (`tsparams.ZtpTestPathNoExtraManifestsCm`) exists.
2. Update the Argo CD clusters app with the non-existent extra manifests configmap reference Git path and wait for synchronization.
3. Pull the `ClusterInstance` CR (`RANConfig.Spoke1Name`) from the hub cluster.
4. Verify that the `ClusterInstance` CR reports a validation failed condition with the message "Validation failed: failed to retrieve ExtraManifest" within `tsparams.SiteconfigOperatorDefaultReconcileTime`.
