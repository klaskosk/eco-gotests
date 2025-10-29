# Test Case Summary for 54296

Test case 54296 is located in `tests/cnf/ran/talm/tests/talm-batching.go` and is named "should continue the CGU when the second batch fails with the Continue batch timeout action".

## Goal

The goal of this test case is to verify that the CGU continues when the second batch fails with the "Continue" batch timeout action, and that the CGU timeout is recalculated for later batches after earlier batches complete.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure the hub and two spokes are present.
- Ensure TALM is at least version 4.12.
- Cleanup test resources on hub and spokes.
- Verify that the temporary namespace (`tsparams.TemporaryNamespace`) does not exist on `Spoke2APIClient`.
- Create the temporary namespace (`tsparams.TemporaryNamespace`) on `Spoke1APIClient` only.

It does not require a git config set up such that X.

## Test Steps

1. Set `expectedTimeout` to 16.
2. Create a CGU with `RANConfig.Spoke1Name` and `RANConfig.Spoke2Name` as clusters, and `tsparams.PolicyName` as the managed policy. Set `MaxConcurrency` to 1, `RemediationStrategy.Timeout` to `expectedTimeout`, and `Enable` to false.
3. Setup the CGU with a catalog source using `helper.SetupCguWithCatSrc`.
4. Wait to enable the CGU using `helper.WaitToEnableCgu`.
5. Wait for the CGU to reach the `tsparams.CguTimeoutReasonCondition` within 21 minutes.
6. Validate that the policy succeeded on spoke1 (catalog source exists on spoke1).
7. Validate that the policy failed on spoke2 (catalog source does not exist on spoke2).
8. Validate that the CGU timeout is recalculated for later batches after earlier batches complete, by checking that the elapsed time between `StartedAt` and `CompletedAt` is approximately equal to `expectedTimeout` minutes, plus or minus `tsparams.TalmDefaultReconcileTime`.
