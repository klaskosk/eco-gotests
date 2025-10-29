# Test Case Summary for 72143

Test case 72143 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts DNS Resolution after Ingress scale-down and scale-up".

## Goal

The goal of this test case is to verify that DNS resolution continues to function correctly from an existing deployment after the SPK Ingress data plane and DNS plane deployments are scaled down to zero replicas and then scaled back up to one replica.

## Test Setup

Prior to the test case, it is assumed that:

- Existing SPK Ingress deployments (`SPKConfig.SPKDataIngressDeployName` and `SPKConfig.SPKDnsIngressDeployName`) are present.
- An existing workload deployment (`SPKConfig.WorkloadDCIDeploymentName`) is present, with pods matching the `wlkdDCILabel` and containing a container named `wlkdDCIContainerName`.

It does not require a git config set up.

## Test Steps

1. The SPK Ingress data plane deployment (`SPKConfig.SPKDataIngressDeployName`) is scaled down to 0 replicas.
2. The SPK Ingress DNS plane deployment (`SPKConfig.SPKDnsIngressDeployName`) is scaled down to 0 replicas.
3. The SPK Ingress data plane deployment is scaled back up to 1 replica.
4. The SPK Ingress DNS plane deployment is scaled back up to 1 replica.
5. The `verifyDNSResolution` function is called to perform DNS lookups and URL access checks from within the existing workload deployment, ensuring that DNS resolution is successful after the Ingress scale-down and scale-up operations.
