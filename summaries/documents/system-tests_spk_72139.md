# Test Case Summary for 72139

Test case 72139 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts DNS Resolution after SPK scale-down and scale-up from existing deployment".

## Goal

The goal of this test case is to verify that DNS resolution continues to function correctly from an existing deployment after the SPK TMM (Traffic Management Microkernel) deployments are scaled down to zero replicas and then scaled back up to one replica.

## Test Setup

Prior to the test case, it is assumed that:

- Existing SPK TMM deployments (`SPKConfig.SPKDataTMMDeployName` and `SPKConfig.SPKDnsTMMDeployName`) are present.
- An existing workload deployment (`SPKConfig.WorkloadDCIDeploymentName`) is present, with pods matching the `wlkdDCILabel` and containing a container named `wlkdDCIContainerName`.

It does not require a git config set up.

## Test Steps

1. The SPK TMM data plane deployment (`SPKConfig.SPKDataTMMDeployName`) is scaled down to 0 replicas.
2. The SPK TMM data plane deployment is scaled back up to 1 replica.
3. The SPK TMM DNS plane deployment (`SPKConfig.SPKDnsTMMDeployName`) is scaled down to 0 replicas.
4. The SPK TMM DNS plane deployment is scaled back up to 1 replica.
5. The `verifyDNSResolution` function is called to perform DNS lookups and URL access checks from within the existing workload deployment, ensuring that DNS resolution is successful after the TMM scale-down and scale-up operations.
