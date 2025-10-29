# Test Case Summary for 71872

Test case 71872 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies all deploymentes are available".

## Goal

The goal of this test case is to verify that all deployments in the cluster are in an "Available" state after an ungraceful reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and the cluster is in a state of recovery where deployments should be re-establishing their availability.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`deployment.WaitForAllDeploymentsAvailable`) to wait for all deployments to become available.
2. It then verifies that no errors occurred during this wait and that all deployments have indeed reached the "Available" state.
3. The intent is to confirm the successful recovery and operational readiness of deployed applications after a disruptive event.
