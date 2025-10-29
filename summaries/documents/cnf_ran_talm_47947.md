# Test Case Summary for 47947

Test case 47947 is located in `tests/cnf/ran/talm/tests/talm-batching.go` and is named "should complete the CGU when two clusters are successful in a single batch".

## Goal

The goal of this test case is to verify that the CGU successfully completes when two clusters are processed in a single batch.

## Test Setup

Prior to the test case, the following changes are needed:

- Ensure the hub and two spokes are present.
- Ensure TALM is at least version 4.12.
- Cleanup test resources on hub and spokes.

It does not require a git config set up such that X.

## Test Steps

1. Create a CGU with `tsparams.PolicyName` as the managed policy. Set `MaxConcurrency` to 1, `RemediationStrategy.Timeout` to 15, and `Enable` to false.
2. Define `policyLabelSelector` with `MatchExpressions` to select clusters with `Key: "common", Operator: "In", Values: []string{"true"}`.
3. Configure `cguBuilder.Definition.Spec.ClusterLabelSelectors` to include two label selectors: one with `MatchLabels` for `RANConfig.Spoke1Name` and another with `MatchExpressions` for `RANConfig.Spoke2Name`.
4. Create a temporary namespace (`tsparams.TemporaryNamespace`) on the hub and create a policy using this namespace definition.
5. Create policy components using `helper.CreatePolicyComponents` with `cguBuilder.Definition.Spec.Clusters` and `policyLabelSelector`.
6. Create the CGU.
7. Wait to enable the CGU using `helper.WaitToEnableCgu`.
8. Wait for the CGU to reach the `tsparams.CguSuccessfulFinishCondition` within 21 minutes, confirming successful completion.
9. Verify the test policy was deleted upon CGU expiration by getting the generated policy name and waiting for it to be deleted.
