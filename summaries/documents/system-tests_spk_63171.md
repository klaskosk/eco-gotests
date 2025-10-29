# Test Case Summary for 63171

Test case 63171 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts DNS resoulution from new deployment".

## Goal

The goal of this case is to verify successful DNS resolution from within a newly created deployment.

## Test Setup

Prior to the test case, the following setup is performed:

- Any existing workload deployment (`SPKConfig.WorkloadDeploymentName`) is deleted.
- A new deployment (`SPKConfig.WorkloadDeploymentName`) is created with one replica, running a container (`wlkdContainerName`) that sleeps indefinitely.

It does not require a git config set up.

## Test Steps

1. The test ensures that no previous workload deployment exists and then creates a new one.
2. It verifies that pods for the new deployment are created and ready.
3. From within one of the newly created pods, a DNS lookup is performed for `SPKConfig.WorkloadTestURL` using `dig`, and the output is asserted to be a valid IPv4 address.
4. An attempt is made to access an external URL (formed from `SPKConfig.WorkloadTestURL` and `SPKConfig.WorkloadTestPort`) using `curl` from within the pod, asserting a 200 or 404 HTTP response code.
