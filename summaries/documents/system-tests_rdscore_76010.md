# Test Case Summary for 76010

Test case 76010 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies workload reachable over BGP route post hard reboot".

## Goal

The goal of this test case is to verify that a workload remains reachable over a BGP route after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that BGP routing (e.g., using FRR) is configured and has recovered.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.ReachURLviaFRRroute`) to validate reachability of a workload via a BGP route.
2. This typically involves attempting to access a service or application that is exposed through a BGP-advertised route, confirming that the BGP routing infrastructure is operational post-reboot.
3. The intent is to confirm the proper functioning of BGP routing for workload access after a disruptive event.
