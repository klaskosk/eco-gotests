# Test Case Summary for 72144

Test case 72144 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts DNS Resolution after Ingress scale-down and scale-up from new deployment".

## Goal

The goal of this test case is to verify that DNS resolution continues to function correctly from a *newly created* deployment after the SPK Ingress data plane and DNS plane deployments are scaled down to zero replicas and then scaled back up to one replica.

## Test Setup

Prior to the test case, the following setup is implicitly handled by the `VerifyDNSResolutionFromNewDeploy` function:

- Any existing workload deployment (`SPKConfig.WorkloadDeploymentName`) is deleted.
- A new deployment (`SPKConfig.WorkloadDeploymentName`) is created with one replica, running a container (`wlkdContainerName`) that sleeps indefinitely.

It does not require a git config set up.

## Test Steps

1. The SPK Ingress data plane deployment (`SPKConfig.SPKDataIngressDeployName`) is scaled down to 0 replicas.
2. The SPK Ingress DNS plane deployment (`SPKConfig.SPKDnsIngressDeployName`) is scaled down to 0 replicas.
3. The SPK Ingress data plane deployment is scaled back up to 1 replica.
4. The SPK Ingress DNS plane deployment is scaled back up to 1 replica.
5. The `VerifyDNSResolutionFromNewDeploy` function is called. This function ensures that no previous workload deployment exists and then creates a new one. It verifies that pods for the new deployment are created and ready.
6. From within one of the newly created pods, a DNS lookup is performed for `SPKConfig.WorkloadTestURL` using `dig`, and the output is asserted to be a valid IPv4 address.
7. An attempt is made to access an external URL (formed from `SPKConfig.WorkloadTestURL` and `SPKConfig.WorkloadTestPort`) using `curl` from within the pod, asserting a 200 or 404 HTTP response code.
