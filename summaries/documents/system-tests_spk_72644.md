# Test Case Summary for 72644

Test case 72644 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts DNS Resolution after SPK TMM pod(s) are deleted from existing deployment".

## Goal

The goal of this test case is to verify that DNS resolution continues to function correctly from an existing deployment after the SPK TMM (Traffic Management Microkernel) data plane and DNS plane pods are deleted.

## Test Setup

Prior to the test case, it is assumed that:

- Existing SPK TMM pods (matching `tmmLabel`) are present in both `SPKConfig.SPKDataNS` and `SPKConfig.SPKDnsNS`.
- An existing workload deployment (`SPKConfig.WorkloadDCIDeploymentName`) is present, with pods matching the `wlkdDCILabel` and containing a container named `wlkdDCIContainerName`.

It does not require a git config set up.

## Test Steps

1. All SPK TMM data plane pods (matching `tmmLabel`) in `SPKConfig.SPKDataNS` are deleted.
2. All SPK TMM DNS plane pods (matching `tmmLabel`) in `SPKConfig.SPKDnsNS` are deleted.
3. The `verifyDNSResolution` function is called to perform DNS lookups and URL access checks from within the existing workload deployment, ensuring that DNS resolution is successful after the TMM pods are deleted.
