# Test Case Summary for 77375

Test case 77375 is located in tests/cnf/ran/oran/tests/oran-post-provision.go and is named "successfully updates ClusterInstance defaults".

## Goal

The goal of this test case is to successfully update the `ClusterInstance` defaults by changing the `TemplateVersion` of the `ProvisioningRequest` and verify that the `ManagedCluster` is updated with the expected label.

## Test Setup

Prior to the test case, the following changes are needed:

- A `ProvisioningRequest` named `tsparams.TestPRName` is pulled and its original spec is saved. It's also verified to be in the `Fulfilled` state.

It does not require a git config set up.

## Test Steps

1. Verify that the test label (`tsparams.TestName`) does not already exist on the `ManagedCluster`.
2. Update the `TemplateVersion` of the `ProvisioningRequest` to `RANConfig.ClusterTemplateAffix + "-" + tsparams.TemplateUpdateDefaults`.
3. Update the `ProvisioningRequest` on the cluster.
4. Wait for the test label (`tsparams.TestName`) to appear on the `ClusterInstance` and then on the `ManagedCluster`.
