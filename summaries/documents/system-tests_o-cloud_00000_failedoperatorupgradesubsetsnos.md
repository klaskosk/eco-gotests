# Test Case Summary for Failed Operator Upgrade in a Subset of SNOs

Test case "Failed operator upgrade in a subset of the SNOs" is located in tests/system-tests/o-cloud/tests/day2-configuration.go and is named "It verifies failed operator upgrade in a subset of the SNOs".

## Goal

The goal of this test case is to verify that the operator upgrade process selectively fails in a subset of SNO clusters when resource limits are intentionally configured to cause failure for specific clusters.

## Test Setup

Prior to the test case, the following changes are needed:

- Downgrade operator images using `downgradeOperatorImages()`.
- Ensure BMHs (OCloudConfig.BmhSpoke1, OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.

It does not require a git config set up.

## Test Steps

1. Verify that both BMHs (OCloudConfig.BmhSpoke1 and OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.
2. Provision two SNO clusters:
    a. Provision the first SNO cluster using `VerifyProvisionSnoCluster`.
    b. Provision the second SNO cluster using `VerifyProvisionSnoCluster`.
3. For both provisioned clusters:
    a. Verify that the OCloud Custom Resources (CRs) exist using `VerifyOcloudCRsExist`.
    b. Verify that the cluster instance creation is completed using `VerifyClusterInstanceCompleted`.
4. Concurrently verify that all policies in the namespaces of both SNO clusters are compliant using `VerifyAllPoliciesInNamespaceAreCompliant`.
5. For both provisioning requests:
    a. Pull the provisioning request using `oran.PullPR`.
    b. Verify that the provisioning request is fulfilled using `VerifyProvisioningRequestIsFulfilled`.
6. Set CPU limits for the PTP operator deployment in `OCloudConfig.ClusterName2` (one of the SNOs) to a low value (1m) to intentionally trigger an upgrade failure for this specific SNO.
7. Create API clients for both SNO clusters.
8. Attempt to upgrade operators in both SNO clusters using `upgradeOperators`.
9. Verify PTP operator versions:
    a. Assert that the PTP operator version in `OCloudConfig.ClusterName1` has been successfully upgraded.
    b. Assert that the PTP operator version in `OCloudConfig.ClusterName2` has not changed, confirming the upgrade failure for the subset.
10. Remove the `tmp/` directory.
11. Concurrently deprovision both SNO clusters using `DeprovisionAiSnoCluster`.
