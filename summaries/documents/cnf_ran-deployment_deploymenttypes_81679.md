# Test Case Summary for 81679

Test case 81679 is located in tests/cnf/ran-deployment/deploymenttypes/tests/deployment-types.go and is named "checks if cluster type is SNOPlusWorker".

## Goal

The goal of this test case is to verify if the cluster type is SNOPlusWorker.

## Test Setup

Prior to the test case, the following changes are needed:

- The test verifies that all policies are compliant for Spoke1 and Spoke2 (if available).
- It retrieves ArgoCD `policiesApp` and `clustersApp` git source details.
- It clones the siteconfig and policy repositories.
- It determines the deployment method and policy template from the cloned repositories.
- It determines the cluster type.

It does not require a git config set up such that X.

## Test Steps

1. Check if `clusterKind` is not empty.
2. If `clusterKind` is not `SNOPlusWorker`, skip the test.
3. Log that the cluster type is `SNOPlusWorker`.
