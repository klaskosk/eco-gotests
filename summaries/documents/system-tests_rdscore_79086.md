# Test Case Summary for 79086

Test case 79086 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies workload reachable over correct BGP route learned by MetalLB FRR post hard reboot".

## Goal

The goal of this test case is to verify that a workload remains reachable over the correct BGP route, as learned by MetalLB FRR, after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that MetalLB FRR is configured for BGP routing and has recovered its routing state.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyMetallbEgressTrafficSegregation`) to validate the egress traffic segregation via MetalLB FRR BGP routes.
2. This typically involves attempting to access a service or application and verifying that the traffic is routed through the expected MetalLB FRR BGP route, confirming the integrity of the routing setup post-reboot.
3. The intent is to confirm the proper functioning of MetalLB FRR BGP routing for workload access and traffic segregation after a disruptive event.
