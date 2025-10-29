# Test Case Summary for 77394

Test case 77394 is located in tests/cnf/ran/oran/tests/oran-provision.go and is named "successfully provisions and generates the correct resources".

## Goal

The goal of this test case is to verify that provisioning with a valid `ProvisioningRequest` successfully provisions and generates the correct resources.

## Test Setup

Prior to the test case, the following changes are needed:

- The test ensures that the `ProvisioningRequest` named `tsparams.TestPRName` exists. If it does not, it will be created.
- If `RANConfig.Spoke1Kubeconfig` or `RANConfig.Spoke1Password` are set, their values are saved to secrets after each test.

It does not require a git config set up.

## Test Steps

1. Pull the `ProvisioningRequest` named `tsparams.TestPRName`. If it does not exist, create a new one using `helper.NewProvisioningRequest` with `tsparams.TemplateValid`.
2. Wait for the `ProvisioningRequest` to be fulfilled.
3. Verify that spoke provisioning succeeded by calling `verifySpokeProvisioning`.
    - This function verifies the creation of the `pull-secret`, `extra-manifests` `ConfigMap`, and the policy `ConfigMap` for spoke 1.
    - It also waits for all policies to achieve a `Compliant` state.
