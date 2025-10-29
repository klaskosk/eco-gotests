# Test Case Summary for 80498

Test case 80498 is located in tests/cnf/ran-deployment/deploymenttypes/tests/deployment-types.go and is named "checks if deployment is multi cluster".

## Goal

The goal of this test case is to verify if the cluster deployment is a multi-cluster deployment.

## Test Setup

Prior to the test case, the following changes are needed:

- The test verifies that all policies are compliant for Spoke1 and Spoke2 (if available).
- It retrieves ArgoCD `policiesApp` and `clustersApp` git source details.
- It clones the siteconfig and policy repositories.
- It determines the deployment method and policy template from the cloned repositories.
- It determines the cluster type.

It does not require a git config set up such that X.

## Test Steps

1. Check if `isMultiCluster` variable is either `SingleCluster` or `MultiCluster`.
2. If `isMultiCluster` is `SingleCluster`, skip the test.
3. Log that the deployment is multi-cluster.
