# Test Case Summary for Failed SNO Provisioning using Image Based Installer

Test case "Verifies the failed provisioning of a single SNO cluster using Image Based Installer" is located in tests/system-tests/o-cloud/tests/sno-provisioning-ibi.go and is named "It verifies the failed provisioning of a single SNO cluster using Image Based Installer".

## Goal

The goal of this test case is to verify that the provisioning of a single SNO cluster using the Image Based Installer (IBI) fails as expected when an invalid configuration is provided.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure BMHs (OCloudConfig.BmhSpoke1, OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.
- If `OCloudConfig.GenerateSeedImage` is true and a base image does not exist, a base image needs to be generated.

It does not require a git config set up.

## Test Steps

1. Verify that both BMHs (OCloudConfig.BmhSpoke1 and OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.
2. If configured, generate the base image by calling `generateBaseImage` if it does not already exist.
3. Install the base image on `OCloudConfig.Spoke2BMC` using `installBaseImage` with `OCloudConfig.IbiBaseImageURL`, `OCloudConfig.VirtualMediaID`, `OCloudConfig.SSHCluster2`, `ocloudparams.SpokeSSHUser`, and `ocloudparams.SpokeSSHPasskeyPath`.
4. Attempt to provision a SNO cluster using `VerifyProvisionSnoCluster` with `OCloudConfig.TemplateName`, `OCloudConfig.TemplateVersionIBIFailure` (indicating an expected failure), `OCloudConfig.NodeClusterName2`, `OCloudConfig.OCloudSiteID`, `ocloudparams.PolicyTemplateParameters`, and `ocloudparams.ClusterInstanceParameters2`.
5. Verify that the OCloud Custom Resources (CRs) for the provisioning request exist using `VerifyOcloudCRsExist`.
6. Verify that the image cluster install is completed using `VerifyImageClusterInstallCompleted`.
7. Verify that the provisioning request times out using `VerifyProvisioningRequestTimeout`.
8. Deprovision the IBI SNO cluster using `DeprovisionIbiSnoCluster` to clean up the failed provisioning attempt.
