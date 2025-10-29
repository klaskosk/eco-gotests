# Test Case Summary for 75378

Test case 75378 is located in tests/cnf/ran/gitopsztp/tests/ztp-siteconfig-negative.go and is named "Validate erroneous/invalid ClusterInstance CR does not block other ClusterInstance CR handling - Verify erroneous/invalid ClusterInstance CR does not block other ClusterInstance CR handling".

## Goal

The goal of this test case is to verify that the siteconfig operator can handle multiple `ClusterInstance` CRs, even if one of them is erroneous or invalid, and that an invalid CR does not prevent the proper handling of other valid `ClusterInstance` CRs.

## Test Setup

Prior to the test case, the original `clustersApp` source is saved and then reset in `BeforeEach` and `AfterEach` blocks respectively. It requires ZTP version 4.17 or later.

It does require a git config set up such that the clusters app can be updated via GitOps.

## Test Steps

1. Check if the ZTP test path for an invalid template reference (`tsparams.ZtpTestPathInvalidTemplateRef`) exists.
2. Update the Argo CD clusters app with the invalid template reference Git path and wait for synchronization.
3. Pull `clusterInstance1` (for `RANConfig.Spoke1Name`) from the hub cluster.
4. Verify that `clusterInstance1` reports a "Validation failed: failed to validate node-level TemplateRef" condition within `tsparams.SiteconfigOperatorDefaultReconcileTime`.
5. Check if the ZTP test path for a valid template reference (`tsparams.ZtpTestPathValidTemplateRef`) exists.
6. Update the Argo CD clusters app with the valid template reference Git path and wait for synchronization.
7. Pull `clusterInstance2` (for `RANConfig.Spoke2Name`) from the hub cluster.
8. Verify that `clusterInstance2` reports a "Provisioning cluster" message with reason "InProgress" within `tsparams.SiteconfigOperatorDefaultReconcileTime`.
