# Test Case Summary for Successful Operator Upgrade

Test case "Successful operator upgrade" is located in tests/system-tests/o-cloud/tests/day2-configuration.go and is named "It verifies successful operator upgrade".

## Goal

The goal of this test case is to verify the successful upgrade of operators on SNO clusters.

## Test Setup

Prior to the test case, the following changes are needed:

- Downgrade operator images using `downgradeOperatorImages()`.
- Ensure BMHs (OCloudConfig.BmhSpoke1, OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.

It does not require a git config set up.

## Test Steps

1. Verify that both BMHs (OCloudConfig.BmhSpoke1 and OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.
2. Provision two SNO clusters:
    a. Provision the first SNO cluster using `VerifyProvisionSnoCluster` with `OCloudConfig.TemplateName`, `OCloudConfig.TemplateVersionDay2`, `OCloudConfig.NodeClusterName1`, `OCloudConfig.OCloudSiteID`, `ocloudparams.PolicyTemplateParameters`, and `ocloudparams.ClusterInstanceParameters1`.
    b. Provision the second SNO cluster using `VerifyProvisionSnoCluster` with `OCloudConfig.TemplateName`, `OCloudConfig.TemplateVersionDay2`, `OCloudConfig.NodeClusterName2`, `OCloudConfig.OCloudSiteID`, `ocloudparams.PolicyTemplateParameters`, and `ocloudparams.ClusterInstanceParameters2`.
3. For both provisioned clusters:
    a. Verify that the OCloud Custom Resources (CRs) exist using `VerifyOcloudCRsExist`.
    b. Verify that the cluster instance creation is completed using `VerifyClusterInstanceCompleted`.
4. Verify that all policies in the namespaces of both SNO clusters are compliant using `VerifyAllPoliciesInNamespaceAreCompliant`.
5. For both provisioning requests:
    a. Pull the provisioning request using `oran.PullPR`.
    b. Verify that the provisioning request is fulfilled using `VerifyProvisioningRequestIsFulfilled`.
