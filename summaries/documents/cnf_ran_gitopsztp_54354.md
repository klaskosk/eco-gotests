# Test Case Summary for 54354

Test case 54354 is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-policies-app.go and is named "Ability to configure local registry via du profile - verifies the image registry exists".

## Goal

The goal of this test case is to verify that a local image registry can be configured via a DU profile using GitOps and that the image registry exists and is available on the spoke cluster.

## Test Setup

Prior to the test case, the original policies app source is saved and then reset after the test. It requires ZTP version 4.10 or later. The test also saves the image registry configuration before modification and restores it afterwards. It also cleans up any image registry leftovers.

It does require a git config set up such that the policies app can be updated via GitOps.

## Test Steps

1. Save the existing image registry configuration.
2. Check if the ZTP test path for image registry exists.
3. Update the Argo CD policies app with the image registry Git path (`tsparams.ZtpTestPathImageRegistry`) and wait for synchronization.
4. Check if the image registry directory (`tsparams.ImageRegistryPath`) is present on spoke 1.
5. Wait for all specified image registry policies (`tsparams.ImageRegistryPolicies`) to exist and become `Compliant`.
6. Wait for the image registry configuration to be `Available` with reason `Ready`.
