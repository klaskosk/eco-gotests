# Test Case Summary for Successful SNO Provisioning using Assisted Installer

Test case "Verifies the successful provisioning of a single SNO cluster using Assisted Installer" is located in tests/system-tests/o-cloud/tests/sno-provisioning-ai.go and is named "It verifies the successful provisioning of a single SNO cluster using Assisted Installer".

## Goal

The goal of this test case is to verify the successful provisioning of a single SNO cluster utilizing the Assisted Installer.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure BMHs (OCloudConfig.BmhSpoke1, OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.

It does not require a git config set up.

## Test Steps

1. Verify that both BMHs (OCloudConfig.BmhSpoke1 and OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.
2. Provision a SNO cluster using `VerifyProvisionSnoCluster` with `OCloudConfig.TemplateName`, `OCloudConfig.TemplateVersionAISuccess`, `OCloudConfig.NodeClusterName1`, `OCloudConfig.OCloudSiteID`, `ocloudparams.PolicyTemplateParameters`, and `ocloudparams.ClusterInstanceParameters1`.
3. Verify that the OCloud Custom Resources (CRs) for the provisioned cluster exist using `VerifyOcloudCRsExist`.
4. Verify that the cluster instance creation is completed using `VerifyClusterInstanceCompleted`.
5. Verify that all policies in the cluster's namespace are compliant using `VerifyAllPoliciesInNamespaceAreCompliant`.
6. Verify that the provisioning request is fulfilled using `VerifyProvisioningRequestIsFulfilled`.
7. Deprovision the AI SNO cluster using `DeprovisionAiSnoCluster` to clean up the provisioned cluster.
