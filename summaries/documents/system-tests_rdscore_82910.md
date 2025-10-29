# Test Case Summary for 82910

Test case 82910 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify Whereabouts Deployment on the same node after node power off".

## Goal

The goal of this test case is to verify inter-pod communication for a Whereabouts-managed deployment on the same node after a node power off event.

## Test Setup

Prior to the test case, this test assumes that Whereabouts is installed and configured, and a deployment managed by Whereabouts is running with all its pods on the same node. The test will simulate a node power off event.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyWhereaboutsInterDeploymentPodCommunicationOnTheSameNodeAfterNodePowerOff` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that pods within a Whereabouts-managed deployment can re-establish and maintain communication on the same node after a node is powered off and restarted.
