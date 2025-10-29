# Test Case Summary for 72141

Test case 72141 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts DNS Resolution with multiple TMM controllers from existing deployment".

## Goal

The goal of this test case is to verify that DNS resolution continues to function correctly from an existing deployment when the SPK TMM (Traffic Management Microkernel) deployments are scaled up to include multiple controllers (2 replicas).

## Test Setup

Prior to the test case, it is assumed that:

- Existing SPK TMM deployments (`SPKConfig.SPKDataTMMDeployName` and `SPKConfig.SPKDnsTMMDeployName`) are present.
- An existing workload deployment (`SPKConfig.WorkloadDCIDeploymentName`) is present, with pods matching the `wlkdDCILabel` and containing a container named `wlkdDCIContainerName`.

It does not require a git config set up.

## Test Steps

1. The SPK TMM data plane deployment (`SPKConfig.SPKDataTMMDeployName`) is scaled up to 2 replicas.
2. The SPK TMM DNS plane deployment (`SPKConfig.SPKDnsTMMDeployName`) is scaled up to 2 replicas.
3. The `verifyDNSResolution` function is called to perform DNS lookups and URL access checks from within the existing workload deployment, ensuring that DNS resolution is successful with multiple TMM controllers.
