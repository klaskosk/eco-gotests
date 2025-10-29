# Test Case Summary for 75374

Test case 75374 is located in tests/cnf/ran/gitopsztp/tests/ztp-cluster-instance-delete.go and is named "Detaching the AI multi-node openshift (MNO) spoke cluster - Validate detaching the AI multi-node openshift spoke cluster".

## Goal

The goal of this test case is to validate the process of detaching an Assisted Installer (AI) multi-node OpenShift (MNO) spoke cluster, ensuring that the relevant AI cluster installation Custom Resources (CRs) are removed while the installed spoke cluster remains functional.

## Test Setup

Prior to the test case, the original clusters app source is saved and then reset after the test. It requires ZTP version 4.17 or later. The test skips if the spoke cluster type is SNO or if the Git path `tsparams.ZtpTestPathDetachAIMNO` does not exist. It also deletes default AI cluster and node level templates ConfigMaps before proceeding.

It does require a git config set up such that the clusters app can be updated via GitOps.

## Test Steps

1. Check the spoke cluster type and skip if it's SNO.
2. Check if the ZTP test path for detaching AI MNO exists.
3. Delete default Assisted Installer cluster level templates ConfigMap CR.
4. Delete default Assisted Installer node level templates ConfigMap CR.
5. Verify that the installed spoke cluster is still functional by getting its OCP version.
6. Update the clusters app Git path with `tsparams.ZtpTestPathDetachAIMNO` and wait for synchronization.
7. Call `validateAISpokeClusterInstallCRsRemoved()` to verify AI spoke cluster install CRs are removed.
8. Get the siteconfig operator pod name and delete it from the `rhacm` namespace on the hub cluster.
9. Wait for 10 seconds to allow the siteconfig operator to reconcile.
10. Check that the default AI cluster and node level templates ConfigMap CRs are recreated successfully.
11. Verify that the installed spoke cluster is still functional.
12. Verify that the spoke cluster namespace CR and cluster instance CR exist on the hub after the siteconfig operator's pod restart.
