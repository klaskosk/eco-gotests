# Test Case Summary for 77374

Test case 77374 is located in tests/cnf/ran/oran/tests/oran-post-provision.go and is named "successfully updates policyTemplateParameters".

## Goal

The goal of this test case is to successfully update the `policyTemplateParameters` of a `ProvisioningRequest` and verify that a `ConfigMap` is updated with the new value.

## Test Setup

Prior to the test case, the following changes are needed:

- A `ProvisioningRequest` named `tsparams.TestPRName` is pulled and its original spec is saved. It's also verified to be in the `Fulfilled` state.
- A `ConfigMap` named `tsparams.TestName` exists and contains `tsparams.TestOriginalValue`.

It does not require a git config set up.

## Test Steps

1. Verify that the test `ConfigMap` (`tsparams.TestName`) exists and has the `tsparams.TestOriginalValue`.
2. Update the `policyTemplateParameters` of the `ProvisioningRequest` to set the `tsparams.TestName` key to `tsparams.TestNewValue`.
3. Update the `ProvisioningRequest` with the modified `policyTemplateParameters`.
4. Wait for the `ProvisioningRequest` to be `Fulfilled` again.
5. Verify that the test `ConfigMap` (`tsparams.TestName`) now has the `tsparams.TestNewValue`.
