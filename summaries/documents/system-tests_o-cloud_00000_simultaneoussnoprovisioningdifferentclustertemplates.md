# Test Case Summary for Simultaneous SNO Provisioning with Different Cluster Templates

Test case "Verifies the successful E2E simultaneous provisioning of SNO clusters with different cluster templates" is located in tests/system-tests/o-cloud/tests/sno-provisioning-ai.go and is named "It verifies the successful E2E simultaneous provisioning of SNO clusters with different cluster templates".

## Goal

The goal of this test case is to verify the successful end-to-end simultaneous provisioning of two SNO clusters, each using a distinct cluster template.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure BMHs (OCloudConfig.BmhSpoke1, OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.

It does not require a git config set up.

## Test Steps

1. Verify that both BMHs (OCloudConfig.BmhSpoke1 and OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.
2. Simultaneously provision two SNO clusters, each with a different cluster template:
    a. Provision the first SNO cluster using `VerifyProvisionSnoCluster` with `OCloudConfig.TemplateName`, `OCloudConfig.TemplateVersionSimultaneous1`, `OCloudConfig.NodeClusterName1`, `OCloudConfig.OCloudSiteID`, `ocloudparams.PolicyTemplateParameters`, and `ocloudparams.ClusterInstanceParameters1`.
    b. Provision the second SNO cluster using `VerifyProvisionSnoCluster` with `OCloudConfig.TemplateName`, `OCloudConfig.TemplateVersionSimultaneous2`, `OCloudConfig.NodeClusterName2`, `OCloudConfig.OCloudSiteID`, `ocloudparams.PolicyTemplateParameters`, and `ocloudparams.ClusterInstanceParameters2`.
3. For both provisioned clusters, verify that the OCloud Custom Resources (CRs) exist using `VerifyOcloudCRsExist`.
4. Concurrently verify that all policies in the namespaces of both SNO clusters are compliant using `VerifyAllPoliciesInNamespaceAreCompliant`.
