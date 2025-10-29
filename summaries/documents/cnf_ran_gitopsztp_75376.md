# Test Case Summary for 75376

Test case 75376 is located in tests/cnf/ran/gitopsztp/tests/ztp-cluster-instance-delete.go and is named "Detaching the AI single-node openshift (SNO) spoke cluster - Validate detaching the AI single-node openshift spoke cluster".

## Goal

The goal of this test case is to validate the process of detaching an Assisted Installer (AI) single-node OpenShift (SNO) spoke cluster, ensuring that the relevant AI cluster installation Custom Resources (CRs) are removed while the installed spoke cluster remains functional.

## Test Setup

Prior to the test case, the original clusters app source is saved and then reset after the test. It requires ZTP version 4.17 or later. The test skips if the spoke cluster type is HighlyAvailableCluster.

It does require a git config set up such that the clusters app can be updated via GitOps.

## Test Steps

1. Check the spoke cluster type and skip if it's a Highly Available Cluster.
2. Check if the ZTP test path for detaching AI SNO exists.
3. Update the clusters app Git path with `tsparams.ZtpTestPathDetachAISNO` and wait for synchronization.
4. Call `validateAISpokeClusterInstallCRsRemoved()` to verify AI spoke cluster install CRs are removed and the spoke cluster is accessible.
