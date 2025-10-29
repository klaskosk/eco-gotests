# Test Case Summary for 72280

Test case 72280 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Assert DNS resolution after Ingress pods were deleted from existing deployment".

## Goal

The goal of this test case is to verify that DNS resolution continues to function correctly from an existing deployment after the SPK Ingress data plane and DNS plane pods are deleted.

## Test Setup

Prior to the test case, it is assumed that:

- Existing SPK Ingress pods (matching `ingressDataLabel` and `ingressDNSLabel`) are present in `SPKConfig.SPKDataNS` and `SPKConfig.SPKDnsNS` respectively.
- An existing workload deployment (`SPKConfig.WorkloadDCIDeploymentName`) is present, with pods matching the `wlkdDCILabel` and containing a container named `wlkdDCIContainerName`.

It does not require a git config set up.

## Test Steps

1. All SPK Ingress data plane pods (matching `ingressDataLabel`) in `SPKConfig.SPKDataNS` are deleted.
2. All SPK Ingress DNS plane pods (matching `ingressDNSLabel`) in `SPKConfig.SPKDnsNS` are deleted.
3. The `verifyDNSResolution` function is called to perform DNS lookups and URL access checks from within the existing workload deployment, ensuring that DNS resolution is successful after the Ingress pods are deleted.
