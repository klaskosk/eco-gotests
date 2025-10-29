# Test Case Summary for Simultaneous SNO Provisioning with Same Cluster Template

Test case "Verifies the successful E2E simultaneous provisioning of SNO clusters with the same cluster template" is located in tests/system-tests/o-cloud/tests/sno-provisioning-ai.go and is named "It verifies the successful E2E simultaneous provisioning of SNO clusters with the same cluster template".

## Goal

The goal of this test case is to verify the successful end-to-end simultaneous provisioning of two SNO clusters using the same cluster template.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure BMHs (OCloudConfig.BmhSpoke1, OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.

It does not require a git config set up.

## Test Steps

1. Verify that both BMHs (OCloudConfig.BmhSpoke1 and OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.
2. Simultaneously provision two SNO clusters:
    a. Provision the first SNO cluster using `VerifyProvisionSnoCluster` with `OCloudConfig.TemplateName`, `OCloudConfig.TemplateVersionAISuccess`, `OCloudConfig.NodeClusterName1`, `OCloudConfig.OCloudSiteID`, `ocloudparams.PolicyTemplateParameters`, and `ocloudparams.ClusterInstanceParameters1`.
    b. Provision the second SNO cluster using `VerifyProvisionSnoCluster` with `OCloudConfig.TemplateName`, `OCloudConfig.TemplateVersionAISuccess`, `OCloudConfig.NodeClusterName2`, `OCloudConfig.OCloudSiteID`, `ocloudparams.PolicyTemplateParameters`, and `ocloudparams.ClusterInstanceParameters2`.
3. For both provisioned clusters:
    a. Verify that the OCloud Custom Resources (CRs) exist using `VerifyOcloudCRsExist`.
4. Concurrently verify that all policies in the namespaces of both SNO clusters are compliant using `VerifyAllPoliciesInNamespaceAreCompliant`.
5. For both provisioning requests, confirm they have been fulfilled using `VerifyProvisioningRequestIsFulfilled`.
