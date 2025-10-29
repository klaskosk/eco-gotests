# Test Case Summary for 77377

Test case 77377 is located in tests/cnf/ran/oran/tests/oran-post-provision.go and is named "successfully adds new manifest to existing PG".

## Goal

The goal of this test case is to successfully add a new manifest to an existing PolicyGenerator (PG) by changing the `TemplateVersion` of the `ProvisioningRequest` and verify that both the original and the new `ConfigMap` exist and have the expected values.

## Test Setup

Prior to the test case, the following changes are needed:

- A `ProvisioningRequest` named `tsparams.TestPRName` is pulled and its original spec is saved. It's also verified to be in the `Fulfilled` state.
- A `ConfigMap` named `tsparams.TestName` exists and contains `tsparams.TestOriginalValue`.
- A second `ConfigMap` named `tsparams.TestName2` does not exist.

It does not require a git config set up.

## Test Steps

1. Verify that the primary test `ConfigMap` (`tsparams.TestName`) exists and has the `tsparams.TestOriginalValue`.
2. Verify that the second test `ConfigMap` (`tsparams.TestName2`) does not exist.
3. Record the current time for `WaitForPhaseAfter`.
4. Update the `TemplateVersion` of the `ProvisioningRequest` to `RANConfig.ClusterTemplateAffix + "-" + tsparams.TemplateAddNew`.
5. Update the `ProvisioningRequest` on the cluster.
6. Wait for the `ProvisioningRequest` to be `Fulfilled` again after the update time.
7. Verify that the primary test `ConfigMap` (`tsparams.TestName`) still has the `tsparams.TestOriginalValue`.
8. Verify that the second test `ConfigMap` (`tsparams.TestName2`) now exists and has the `tsparams.TestOriginalValue`.
