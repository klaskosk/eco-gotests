# Test Case Summary for 77376

Test case 77376 is located in tests/cnf/ran/oran/tests/oran-post-provision.go and is named "successfully updates existing PG manifest".

## Goal

The goal of this test case is to successfully update an existing PolicyGenerator (PG) manifest by changing the `TemplateVersion` of the `ProvisioningRequest` and verify that the associated `ConfigMap` reflects the new value.

## Test Setup

Prior to the test case, the following changes are needed:

- A `ProvisioningRequest` named `tsparams.TestPRName` is pulled and its original spec is saved. It's also verified to be in the `Fulfilled` state.
- A `ConfigMap` named `tsparams.TestName` exists and contains `tsparams.TestOriginalValue`.

It does not require a git config set up.

## Test Steps

1. Verify that the test `ConfigMap` (`tsparams.TestName`) exists and has the `tsparams.TestOriginalValue`.
2. Record the current time for `WaitForPhaseAfter`.
3. Update the `TemplateVersion` of the `ProvisioningRequest` to `RANConfig.ClusterTemplateAffix + "-" + tsparams.TemplateUpdateExisting`.
4. Update the `ProvisioningRequest` on the cluster.
5. Wait for the `ProvisioningRequest` to be `Fulfilled` again after the update time.
6. Verify that the test `ConfigMap` (`tsparams.TestName`) now has the `tsparams.TestNewValue`.
