# Test Case Summary for 74753

Test case 74753 is located in `tests/cnf/ran/talm/tests/talm-batching.go` and is named "should continue the CGU when the first batch fails with the Continue batch timeout action".

## Goal

The goal of this test case is to verify that when the first batch fails with the "Continue" batch timeout action, the CGU continues to process subsequent batches, and the second batch completes successfully.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure the hub and two spokes are present.
- Ensure TALM is at least version 4.12.
- Cleanup test resources on hub and spokes.
- Verify that the temporary namespace (`tsparams.TemporaryNamespace`) does not exist on `Spoke1APIClient`.
- Create the temporary namespace (`tsparams.TemporaryNamespace`) on `Spoke2APIClient` only.

It does not require a git config set up such that X.

## Test Steps

1. Create a CGU with `RANConfig.Spoke1Name` and `RANConfig.Spoke2Name` as clusters, and `tsparams.PolicyName` as the managed policy. Set `MaxConcurrency` to 1, `RemediationStrategy.Timeout` to 9, `Enable` to false.
2. Setup the CGU with a catalog source using `helper.SetupCguWithCatSrc`.
3. Wait to enable the CGU using `helper.WaitToEnableCgu`.
4. Wait for the CGU to reach the `tsparams.CguTimeoutReasonCondition` within 16 minutes.
5. Validate that the policy succeeded on spoke2 (catalog source exists on spoke2).
6. Validate that the policy failed on spoke1 (catalog source does not exist on spoke1).
