# Test Case Summary for 75342

Test case 75342 is located in tests/cnf/ran/gitopsztp/tests/ztp-siteconfig-day-two.go and is named "Verify modification of cluster labels in ClusterInstance CR using git flows after installation - Verify modification of cluster labels in ClusterInstance CR using git flows after installation".

## Goal

The goal of this test case is to verify that cluster labels in a `ClusterInstance` CR can be modified using Git flows after installation, and these modifications are reflected in both the `ClusterInstance` CR and the `ManagedCluster` CR.

## Test Setup

Prior to the test case, the original `clustersApp` source is saved and reset in `BeforeEach` and `AfterEach` blocks respectively. It requires ZTP version 4.17 or later. The `AfterEach` block also verifies that the newly added cluster label is removed from both the `ClusterInstance` CR and the `ManagedCluster` CR.

It does require a git config set up such that the clusters app can be updated via GitOps.

## Test Steps

1. Check if the ZTP test path for adding a new cluster label (`tsparams.ZtpTestPathNewClusterLabel`) exists.
2. Update the Argo CD clusters app with the Git path referencing the new custom label addition (`tsparams.ZtpTestPathNewClusterLabel`) and wait for synchronization.
3. Check that the `ClusterInstance` CR is updated with the newly added cluster label (`tsparams.TestLabelKey`) on the hub cluster.
4. Check that the `ManagedCluster` CR is updated with the newly added spoke cluster label (`tsparams.TestLabelKey`) on the hub cluster.
