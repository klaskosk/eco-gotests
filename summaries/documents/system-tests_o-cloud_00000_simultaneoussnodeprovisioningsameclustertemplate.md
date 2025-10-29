# Test Case Summary for Simultaneous SNO Deprovisioning with Same Cluster Template

Test case "Verifies the successful E2E simultaneous deprovisioning of SNO clusters with the same cluster template" is located in tests/system-tests/o-cloud/tests/sno-provisioning-ai.go and is named "It verifies the successful E2E simultaneous deprovisioning of SNO clusters with the same cluster template".

## Goal

The goal of this test case is to verify the successful end-to-end simultaneous deprovisioning of two SNO clusters that were provisioned using the same cluster template.

## Test Setup

This test case assumes two SNO clusters have been successfully provisioned with the same cluster template and their provisioning requests are fulfilled.

It does not require a git config set up.

## Test Steps

1. Retrieve `provisioningRequest1` and `provisioningRequest2` using `oran.PullPR` and verify that both provisioning requests are fulfilled.
2. Verify that both `provisioningRequest1` and `provisioningRequest2` are using the same cluster template name and version.
3. For both provisioning requests:
    a. Retrieve the Node Allocation Request, allocated Nodes, and namespace using `VerifyOcloudCRsExist`.
    b. Retrieve the Cluster Instance using `VerifyClusterInstanceCompleted`.
    c. Retrieve the BMHs from the allocated nodes using `GetBMHsFromAllocatedNodes`.
4. Concurrently verify the deletion or appropriate state of resources for both clusters:
    a. Verify that both provisioning requests are deleted using `VerifyProvisioningRequestIsDeleted`.
    b. Verify that both namespaces no longer exist using `VerifyNamespaceDoesNotExist`.
    c. Verify that both cluster instances no longer exist using `VerifyClusterInstanceDoesNotExist`.
    d. Verify that both node allocation requests no longer exist using `VerifyNodeAllocationRequestDoesNotExist`.
    e. Verify that both BMHs are in the `Deprovisioned` state using `VerifyBmhIsInDeprovisionedState`.
