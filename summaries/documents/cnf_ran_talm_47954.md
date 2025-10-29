# Test Case Summary for 47954

Test case 47954 is located in `tests/cnf/ran/talm/tests/talm-batching.go` and is named "should report the timeout value when one cluster is in a batch and it times out".

## Goal

The goal of this test case is to verify that the CGU reports the timeout value when a single cluster in a batch times out, and that the total runtime is approximately equal to the expected timeout.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure the hub and two spokes are present.
- Ensure TALM is at least version 4.12.
- Cleanup test resources on hub and spokes.
- Verify that the temporary namespace (`tsparams.TemporaryNamespace`) does not exist on `Spoke1APIClient`.

It does not require a git config set up such that X.

## Test Steps

1. Set `expectedTimeout` to 8.
2. Create a CGU with `RANConfig.Spoke1Name` as the cluster and `tsparams.PolicyName` as the managed policy. Set `MaxConcurrency` to 1, `RemediationStrategy.Timeout` to `expectedTimeout`, and `Enable` to false.
3. Setup the CGU with a catalog source using `helper.SetupCguWithCatSrc`.
4. Wait to enable the CGU using `helper.WaitToEnableCgu`.
5. Wait for the CGU to reach the `tsparams.CguTimeoutReasonCondition` within 11 minutes.
6. Validate that the timeout occurred after just the first reconcile, by checking that the elapsed time between `StartedAt` and `CompletedAt` is approximately equal to `expectedTimeout` minutes, plus or minus `tsparams.TalmDefaultReconcileTime`.
7. Verify the test policy was deleted upon CGU expiration by getting the generated policy name and waiting for it to be deleted.
