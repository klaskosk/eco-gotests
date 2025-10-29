# Test Case Summary for 77373

Test case 77373 is located in tests/cnf/ran/oran/tests/oran-post-provision.go and is named "successfully updates clusterInstanceParameters".

## Goal

The goal of this test case is to successfully update the `clusterInstanceParameters` of a `ProvisioningRequest` and verify that the associated `ManagedCluster` reflects the changes by having the correct label.

## Test Setup

Prior to the test case, the following changes are needed:

- A `ProvisioningRequest` named `tsparams.TestPRName` is pulled and its original spec is saved. It's also verified to be in the `Fulfilled` state.

It does not require a git config set up.

## Test Steps

1. Verify that the test label (`tsparams.TestName`) does not already exist on the `ManagedCluster`.
2. Retrieve the `TemplateParameters` from the `ProvisioningRequest`.
3. Update the `extraLabels` within the `clusterInstanceParameters` to include the `ManagedCluster` label with `tsparams.TestName`.
4. Update the `ProvisioningRequest` with the modified `TemplateParameters`.
5. Wait for the `ClusterInstance` and `ManagedCluster` to have the `tsparams.TestName` label.
