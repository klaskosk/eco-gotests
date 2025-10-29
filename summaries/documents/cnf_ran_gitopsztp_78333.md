# Test Case Summary for 78333

Test case 78333 is located in tests/cnf/ran/gitopsztp/tests/IBBF-e2e-test.go and is named "tests HW replacement via IBBF".

## Goal

The goal of this test case is to test the hardware replacement scenario using the Image-Based Break/Fix (IBBF) flow, specifically verifying that a preserved ConfigMap retains its data but gets a new timestamp, and that the cluster identity (ClusterID, InfraID, and UID) remains consistent after reinstallation.

## Test Setup

Prior to the test case, the clusters app is pulled and the existence of the Git path (`tsparams.ZtpTestPathIBBFe2e`) is checked. The `allowReinstalls` flag in the `siteconfig-operator-configuration` ConfigMap is set to "true". A test ConfigMap (`tsparams.TestCMName`) with a `siteconfig.open-cluster-management.io/preserve` label and data `testValue: true` is created, and its creation timestamp is recorded. The cluster deployment's original Cluster ID, Infra ID, and UID are also retrieved. After each test, the created test ConfigMap is deleted.

It does require a git config set up such that the clusters app can be updated via GitOps.

## Test Steps

1. Enable cluster reinstallation in `SiteconfigOperator` by updating the `siteconfig-operator-configuration` ConfigMap.
2. Create a test ConfigMap (`tsparams.TestCMName`) with a `siteconfig.open-cluster-management.io/preserve` label.
3. Get the cluster deployment identity (ClusterID, InfraID, UID) before IBBF.
4. Change the clusters app to point to the IBBF test target directory (`tsparams.ZtpTestPathIBBFe2e`) and wait for synchronization.
5. Wait for the `ClusterInstance` to trigger re-installation (condition `ReinstallRequestProcessed` with reason `Completed`).
6. Wait for the `ClusterInstance` to start provisioning (condition `ClusterProvisioned` with reason `InProgress`).
7. Wait for the `ClusterInstance` to finish provisioning (condition `ClusterProvisioned` with reason `Completed`).
8. Verify that the test ConfigMap was preserved post-IBBF, ensuring its data is intact but its `CreationTimestamp` is different from the original.
9. Compare the cluster identity (ClusterID, InfraID, UID) post-IBBF, asserting that they remain the same as the original values.
