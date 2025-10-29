# Test Case Summary for 76011

Test case 76011 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies workload reachable over BGP route post graceful reboot".

## Goal

The goal of this test case is to verify that a workload remains reachable over a BGP route after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that BGP routing (e.g., using FRR) is configured and has recovered.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.ReachURLviaFRRroute`) to validate reachability of a workload via a BGP route.
2. This typically involves deploying an application that uses a BGP route for external connectivity.
3. The test then verifies that this application can successfully reach external resources via the configured BGP route after a graceful reboot, confirming its resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable external connectivity via BGP routes after a graceful reboot.
