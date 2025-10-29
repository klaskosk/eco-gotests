# Test Case Summary for 72645

Test case 72645 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts DNS Resolution after SPK TMM pod(s) are deleted from new deployment".

## Goal

The goal of this test case is to verify that DNS resolution continues to function correctly from a *newly created* deployment after the SPK TMM (Traffic Management Microkernel) data plane and DNS plane pods are deleted.

## Test Setup

Prior to the test case, it is implicitly handled by the `VerifyDNSResolutionFromNewDeploy` function that a new workload deployment (`SPKConfig.WorkloadDeploymentName`) is created with one replica.

It does not require a git config set up.

## Test Steps

1. All SPK TMM data plane pods (matching `tmmLabel`) in `SPKConfig.SPKDataNS` are deleted.
2. All SPK TMM DNS plane pods (matching `tmmLabel`) in `SPKConfig.SPKDnsNS` are deleted.
3. The `VerifyDNSResolutionFromNewDeploy` function is called. This function ensures that no previous workload deployment exists and then creates a new one. It verifies that pods for the new deployment are created and ready.
4. From within one of the newly created pods, a DNS lookup is performed for `SPKConfig.WorkloadTestURL` using `dig`, and the output is asserted to be a valid IPv4 address.
5. An attempt is made to access an external URL (formed from `SPKConfig.WorkloadTestURL` and `SPKConfig.WorkloadTestPort`) using `curl` from within the pod, asserting a 200 or 404 HTTP response code.
