# Test Case Summary for 82736

Test case 82736 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies connectivity between pods from deployment scheduled on the same node post graceful reboot".

## Goal

The goal of this test case is to verify network connectivity between pods belonging to a deployment, scheduled on the same node, after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that a deployment with pods deployed on the same node is present and has recovered its state.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyWhereaboutsInterDeploymentPodCommunicationOnTheSameNodeAfterNodePowerOff`) to validate the connectivity between the deployment pods on the same node.
2. This typically involves performing network connectivity tests (e.g., ping, curl) between the pods within the deployment.
3. The intent is to confirm that the deployment pods can re-establish communication and maintain their network identity on the same node after a controlled cluster restart.
