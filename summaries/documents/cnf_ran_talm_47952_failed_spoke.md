# Test Case Summary for 47952 (Failed Spoke)

Test case 47952 is located in `tests/cnf/ran/talm/tests/talm-batching.go` and is named "should report the failed spoke when one spoke in a batch times out".

## Goal

The goal of this test case is to verify that when one spoke in a batch times out, the CGU reports the failed spoke, but the other spoke in the same batch completes successfully.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure the hub and two spokes are present.
- Ensure TALM is at least version 4.12.
- Cleanup test resources on hub and spokes.
- Verify that the temporary namespace (`tsparams.TemporaryNamespace`) does not exist on `Spoke2APIClient`.
- Create the temporary namespace (`tsparams.TemporaryNamespace`) on `Spoke1APIClient` only.

It does not require a git config set up such that X.

## Test Steps

1. Create a CGU with `RANConfig.Spoke1Name` and `RANConfig.Spoke2Name` as clusters, and `tsparams.PolicyName` as the managed policy. Set `MaxConcurrency` to 2, `RemediationStrategy.Timeout` to 9, and `Enable` to false.
2. Setup the CGU with a catalog source using `helper.SetupCguWithCatSrc`.
3. Wait to enable the CGU using `helper.WaitToEnableCgu`.
4. Wait for the CGU to reach the `tsparams.CguTimeoutReasonCondition` within 16 minutes.
5. Validate that the policy succeeded on spoke1 (catalog source exists on spoke1).
6. Validate that the policy failed on spoke2 (catalog source does not exist on spoke2).
