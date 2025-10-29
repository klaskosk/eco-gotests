# Test Case Summary for 66091

Test case 66091 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts DNS resolution from existing deployment".

## Goal

The goal of this test case is to verify successful DNS resolution from within an existing deployment.

## Test Setup

Prior to the test case, it is assumed that a deployment named `SPKConfig.WorkloadDCIDeploymentName` already exists in `SPKConfig.Namespace`, with pods matching the `wlkdDCILabel` and containing a container named `wlkdDCIContainerName`.

It does not require a git config set up.

## Test Steps

1. The test first asserts that the existing deployment `SPKConfig.WorkloadDCIDeploymentName` is present.
2. It then finds pods associated with this deployment based on the `wlkdDCILabel`.
3. From within one of these existing pods, a DNS lookup is performed for `SPKConfig.WorkloadTestURL` using `dig`, and the output is asserted to be a valid IPv4 address.
4. An attempt is made to access an external URL (formed from `SPKConfig.WorkloadTestURL` and `SPKConfig.WorkloadTestPort`) using `curl` from within the pod, asserting a 200 or 404 HTTP response code.
