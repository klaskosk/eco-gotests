# Test Case Summary for 47955

Test case 47955 is located in `tests/cnf/ran/talm/tests/talm-batching.go` and is named "should report the missing policy".

## Goal

The goal of this test case is to verify that the CGU reports a missing policy when a managed policy that does not exist is specified.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure the hub and two spokes are present.
- Ensure TALM is at least version 4.12.
- Cleanup test resources on hub and spokes.

It does not require a git config set up such that X.

## Test Steps

1. Create and enable a CGU with `RANConfig.Spoke1Name` as the cluster and "non-existent-policy" as the managed policy. Set the `RemediationStrategy.Timeout` to 1.
2. Wait for the CGU status to report the `tsparams.CguNonExistentPolicyCondition` within 2 minutes and verify that no error occurred, which confirms the CGU reports the missing policy.
