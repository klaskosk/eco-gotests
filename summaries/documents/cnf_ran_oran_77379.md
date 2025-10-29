# Test Case Summary for 77379

Test case 77379 is located in tests/cnf/ran/oran/tests/oran-post-provision.go and is named "successfully rolls back failed ProvisioningRequest update".

## Goal

The goal of this test case is to successfully roll back a failed `ProvisioningRequest` update. This is achieved by attempting an update that causes a policy to become non-compliant due to an immutable field, and then verifying that the `ProvisioningRequest` is automatically restored to its original state.

## Test Setup

Prior to the test case, the following changes are needed:

- A `ProvisioningRequest` named `tsparams.TestPRName` is pulled and its original spec is saved. It's also verified to be in the `Fulfilled` state.

It does not require a git config set up.

## Test Steps

1. Update the `policyTemplateParameters` of the `ProvisioningRequest` by setting `tsparams.HugePagesSizeKey` to "2G". This change is designed to cause a policy to become non-compliant due to an immutable field.
2. Update the `ProvisioningRequest` on the cluster.
3. Wait for a policy on `RANConfig.Spoke1Name` to transition to `NonCompliant` status, specifically verifying that the non-compliance is due to an immutable field using `helper.WaitForNoncompliantImmutable`.
4. The test cleanup (`AfterEach` block) automatically restores the `ProvisioningRequest` to its original state, demonstrating the successful rollback.
