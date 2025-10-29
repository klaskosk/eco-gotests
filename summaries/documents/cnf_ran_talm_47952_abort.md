# Test Case Summary for 47952 (Abort)

Test case 47952 is located in `tests/cnf/ran/talm/tests/talm-batching.go` and is named "should abort CGU when the first batch fails with the Abort batch timeout action".

## Goal

The goal of this test case is to verify that the CGU aborts when the first batch fails with the "Abort" batch timeout action, and that the total runtime is approximately equal to one reconcile loop.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure the hub and two spokes are present.
- Ensure TALM is at least version 4.12.
- Cleanup test resources on hub and spokes.
- Verify that the temporary namespace (`tsparams.TemporaryNamespace`) does not exist on `Spoke1APIClient`.
- Create the temporary namespace (`tsparams.TemporaryNamespace`) on `Spoke2APIClient` only.

It does not require a git config set up such that X.

## Test Steps

1. Create a CGU with `RANConfig.Spoke1Name` and `RANConfig.Spoke2Name` as clusters, and `tsparams.PolicyName` as the managed policy. Set `MaxConcurrency` to 1, `RemediationStrategy.Timeout` to 9, `Enable` to false, and `BatchTimeoutAction` to "Abort".
2. Setup the CGU with a catalog source using `helper.SetupCguWithCatSrc`.
3. Wait to enable the CGU using `helper.WaitToEnableCgu`.
4. Wait for the CGU to reach the `tsparams.CguTimeoutReasonCondition` within 11 minutes.
5. Validate that the policy failed on spoke1 (catalog source does not exist on spoke1) and on spoke2 (catalog source does not exist on spoke2).
6. Validate that the timeout occurred after just the first reconcile, by checking that the elapsed time between `StartedAt` and `CompletedAt` is approximately equal to `tsparams.TalmDefaultReconcileTime`.
7. Validate that the timeout message matched the abort message by waiting for `tsparams.CguTimeoutMessageCondition`.
