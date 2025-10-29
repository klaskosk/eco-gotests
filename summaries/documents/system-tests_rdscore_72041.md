# Test Case Summary for 72041

Test case 72041 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies all deploymentes are available".

## Goal

The goal of this test case is to verify that all deployments in the cluster are in an "Available" state after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and the cluster is in a state of recovery where deployments should be re-establishing their availability.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.WaitAllDeploymentsAreAvailable`) to wait for all deployments to become available.
2. This typically involves monitoring the status of all `Deployment` resources in the cluster until they report an "Available" status within a defined timeout.
3. The intent is to confirm the proper recovery and availability of all deployed applications after a controlled cluster restart.
