# Test Case Summary for 77393

Test case 77393 is located in tests/cnf/ran/oran/tests/oran-provision.go and is named "recovers provisioning when invalid ProvisioningRequest is updated".

## Goal

The goal of this test case is to verify that the system recovers provisioning when an invalid ProvisioningRequest is updated with valid parameters.

## Test Setup

Prior to the test case, the following changes are needed:

- The ProvisioningRequest named `tsparams.TestPRName` must not exist.

It does not require a git config set up.

## Test Steps

1. Verify the `ProvisioningRequest` named `tsparams.TestPRName` does not already exist. If it does, skip the test.
2. Create a `ProvisioningRequest` with an invalid `policyTemplateParameters` by providing an integer value where a string is expected.
3. Wait for the `ProvisioningRequest` to transition to the `StateFailed` phase.
4. Update the `ProvisioningRequest` with valid `policyTemplateParameters`.
5. Wait for the `ProvisioningRequest` to start progressing, transitioning to the `StateProgressing` phase.
