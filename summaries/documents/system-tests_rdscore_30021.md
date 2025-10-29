# Test Case Summary for 30021

Test case 30021 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies graceful cluster reboot".

## Goal

The goal of this test case is to verify that the cluster can perform a graceful reboot and recover to a functional state.

## Test Setup

Prior to the test case, this test assumes a running cluster configured for graceful reboots.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifySoftReboot`) to initiate and verify a graceful cluster reboot.
2. This typically involves performing a controlled restart of cluster nodes, allowing for services to gracefully shut down and restart.
3. The test then verifies the overall health and functionality of the cluster post-reboot, ensuring all components are operational and no data loss or service disruption occurs.
4. The intent is to confirm the reliability and recovery capabilities of the cluster during a planned, controlled reboot.
