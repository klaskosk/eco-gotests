# Test Case Summary for 47949

Test case 47949 is located in `tests/cnf/ran/talm/tests/talm-batching.go` and is named "should report a missing spoke".

## Goal

The goal of this test case is to verify that the CGU reports a missing spoke when a non-existent cluster is specified in the CGU definition, even when a managed policy is also non-existent.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure the hub and two spokes are present.
- Ensure TALM is at least version 4.12.
- Cleanup test resources on hub and spokes.

It does not require a git config set up such that X.

## Test Steps

1. Create a CGU with a non-existent cluster (`tsparams.NonExistentClusterName`) and a non-existent managed policy (`tsparams.NonExistentPolicyName`). Set the `RemediationStrategy.Timeout` to 1.
2. Wait for the CGU to reach the `tsparams.CguNonExistentClusterCondition` within 3 times the default TALM reconcile time and verify that no error occurred, which confirms that the CGU reports the missing spoke.
