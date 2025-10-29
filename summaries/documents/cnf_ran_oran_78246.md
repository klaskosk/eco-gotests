# Test Case Summary for 78246

Test case 78246 is located in `tests/cnf/ran/oran/tests/oran-pre-provision.go` and is named "successfully generates ClusterInstance provisioning without HardwareTemplate".

## Goal

The goal of this test case is to verify that a `ClusterInstance` is successfully generated even when provisioning without a `HardwareTemplate`.

## Test Setup

Prior to the test case, an O2IMS API client is created and configured. This client is used to interact with the O2IMS API and create a `ProvisioningClient`.

After the test case, a `ProvisioningRequest` with `tsparams.TestPRName` is pulled using `oran.PullPR`. If it exists, it is deleted using `prBuilder.DeleteAndWait`, which deletes the resource and waits for its removal from the cluster.

It does not require a git config set up.

## Test Steps

1. A `ProvisioningRequest` is created using `helper.NewNoTemplatePR` with `tsparams.TemplateNoHWTemplate` as the template version. This helper function creates a `ProvisioningRequestBuilder` and sets template parameters with intentionally incorrect BMC and network data to ensure a `ClusterInstance` is generated but not actually provisioned.
2. The `Create()` method is called on the `ProvisioningRequestBuilder` to create the resource on the cluster. The test asserts that no error occurred during this creation.
3. The test then waits for the associated `ClusterInstance` to be created and validated using `helper.WaitForValidPRClusterInstance`. This function polls the `ClusterInstance` until it has the `RenderedTemplatesApplied` condition set to `ConditionTrue`. The test asserts that no error occurred during this wait, confirming the successful generation and validation of the `ClusterInstance`.
