# Test Case Summary for 82714

Test case 82714 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify Whereabouts Deployment on the same node".

## Goal

The goal of this test case is to verify inter-pod communication for a Whereabouts-managed deployment when all pods are scheduled on the same node.

## Test Setup

Prior to the test case, this test assumes that Whereabouts is installed and configured, and a deployment managed by Whereabouts is running with all its pods on the same node.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyWhereaboutsInterDeploymentPodCommunicationOnTheSameNode` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that pods within a Whereabouts-managed deployment can communicate with each other when co-located on the same node.
