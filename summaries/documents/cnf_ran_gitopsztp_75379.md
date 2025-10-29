# Test Case Summary for 75379

Test case 75379 is located in tests/cnf/ran/gitopsztp/tests/ztp-siteconfig-negative.go and is named "Validate two ClusterInstance CR with duplicate cluster name - Verify two ClusterInstance CR with duplicate cluster name".

## Goal

The goal of this test case is to verify that the siteconfig operator correctly handles the creation of two `ClusterInstance` CRs with duplicate cluster names, ensuring that the second, duplicate CR fails dry-run validation and reports an appropriate error message.

## Test Setup

Prior to the test case, the original `clustersApp` source is saved and then reset in `BeforeEach` and `AfterEach` blocks respectively. It requires ZTP version 4.17 or later.

It does require a git config set up such that the clusters app can be updated via GitOps.

## Test Steps

1. Check if the ZTP test path for a unique cluster name (`tsparams.ZtpTestPathUniqueClusterName`) exists.
2. Update the Argo CD clusters app with the unique cluster name Git path and wait for synchronization.
3. Pull `clusterInstance1` (for `RANConfig.Spoke1Name`) from the hub cluster.
4. Verify that `clusterInstance1` reports a "Provisioning cluster" message with reason "InProgress" within `tsparams.SiteconfigOperatorDefaultReconcileTime`.
5. Check if the ZTP test path for a duplicate cluster name (`tsparams.ZtpTestPathDuplicateClusterName`) exists.
6. Update the Argo CD clusters app with the duplicate cluster name Git path and wait for synchronization.
7. Pull `clusterInstance2` (for `RANConfig.Spoke2Name`) from the hub cluster.
8. Verify that `clusterInstance2` reports a "Rendered manifests failed dry-run validation" message with reason "Failed" within `tsparams.SiteconfigOperatorDefaultReconcileTime`.
