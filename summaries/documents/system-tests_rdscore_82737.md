# Test Case Summary for 82737

Test case 82737 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies connectivity between pods from deployment scheduled on different nodes post graceful reboot".

## Goal

The goal of this test case is to verify network connectivity between pods belonging to a deployment, scheduled on different nodes, after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that a deployment with pods deployed on different nodes is present and has recovered its state.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyWhereaboutsInterDeploymentPodCommunicationOnDifferentNodesAfterNodePowerOff`) to validate the connectivity between the deployment pods across different nodes.
2. This typically involves performing network connectivity tests (e.g., ping, curl) between the pods within the deployment that are distributed on different nodes.
3. The intent is to confirm that the deployment pods can re-establish communication and maintain their network identity across different nodes after a controlled cluster restart.
