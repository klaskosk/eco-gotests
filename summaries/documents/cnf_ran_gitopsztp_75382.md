# Test Case Summary for 75382

Test case 75382 is located in tests/cnf/ran/gitopsztp/tests/ztp-siteconfig-failover.go and is named "Verify siteconfig operator's recovery mechanism by referencing non-existent cluster template configmap CR".

## Goal

The goal of this test case is to verify that the siteconfig operator's recovery mechanism works correctly when a `ClusterInstance` CR references a non-existent cluster template configmap, and that a proper validation error is reported.

## Test Setup

Prior to the test case, the original `clustersApp` source is saved and then reset in `BeforeEach` and `AfterEach` blocks respectively. It requires ZTP version 4.17 or later.

It does require a git config set up such that the clusters app can be updated via GitOps.

## Test Steps

1. Check if the non-existent cluster template configmap reference Git path (`tsparams.ZtpTestPathNoClusterTemplateCm`) exists.
2. Update the Argo CD clusters app with the non-existent cluster template configmap reference Git path and wait for synchronization.
3. Pull the `ClusterInstance` CR (`RANConfig.Spoke1Name`) from the hub cluster.
4. Verify that the `ClusterInstance` CR reports a validation failed condition with the message "Validation failed: failed to validate cluster-level TemplateRef" within `tsparams.SiteconfigOperatorDefaultReconcileTime`.
