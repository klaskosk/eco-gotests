# Test Case Summary for Failed SNO Provisioning using Assisted Installer

Test case "Verifies the failed provisioning of a single SNO cluster using Assisted Installer" is located in tests/system-tests/o-cloud/tests/sno-provisioning-ai.go and is named "It verifies the failed provisioning of a single SNO cluster using Assisted Installer".

## Goal

The goal of this test case is to verify that the provisioning of a single SNO cluster using Assisted Installer fails as expected when an invalid configuration is provided.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure BMHs (OCloudConfig.BmhSpoke1, OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.

It does not require a git config set up.

## Test Steps

1. Verify that both BMHs (OCloudConfig.BmhSpoke1 and OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.
2. Attempt to provision a SNO cluster using `VerifyProvisionSnoCluster` with `OCloudConfig.TemplateName`, `OCloudConfig.TemplateVersionAIFail` (indicating an expected failure), `OCloudConfig.NodeClusterName1`, `OCloudConfig.OCloudSiteID`, `ocloudparams.PolicyTemplateParameters`, and `ocloudparams.ClusterInstanceParameters1`.
3. Verify that the OCloud Custom Resources (CRs) for the provisioning request exist using `VerifyOcloudCRsExist`.
4. Verify that the provisioning request times out using `VerifyProvisioningRequestTimeout`.
5. Deprovision the AI SNO cluster using `DeprovisionAiSnoCluster` to clean up the failed provisioning attempt.
