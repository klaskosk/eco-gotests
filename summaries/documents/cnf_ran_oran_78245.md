# Test Case Summary for 78245

Test case 78245 is located in `tests/cnf/ran/oran/tests/oran-pre-provision.go` and is named "fails to provision without a HardwareTemplate when required schema is missing".

## Goal

The goal of this case is to verify that provisioning fails when a `HardwareTemplate` is missing and a required schema is not provided in the `ClusterTemplate`.

## Test Setup

Prior to the test case, an O2IMS API client is created and configured. This client is used to interact with the O2IMS API and create a `ProvisioningClient`.

It does not require a git config set up.

## Test Steps

1. The test constructs the `clusterTemplateName` and `clusterTemplateNamespace` using `tsparams.ClusterTemplateName`, `RANConfig.ClusterTemplateAffix`, and `tsparams.TemplateMissingSchema`.
2. It then pulls the `ClusterTemplate` using `oran.PullClusterTemplate` with the constructed name and namespace. This function retrieves an existing `ClusterTemplate` from the cluster.
3. The test waits for the `ClusterTemplate` to have the `tsparams.CTInvalidSchemaCondition` using `clusterTemplate.WaitForCondition`. This function polls the `ClusterTemplate` until the specified condition is met or a timeout occurs.
4. The test asserts that no error occurred during the wait, confirming that the `ClusterTemplate` validation failed due to an invalid schema.
