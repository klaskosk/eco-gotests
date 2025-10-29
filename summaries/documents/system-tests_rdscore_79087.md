# Test Case Summary for 79087

Test case 79087 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies workload reachable over correct BGP route learned by MetalLB FRR post graceful reboot".

## Goal

The goal of this test case is to verify that a workload remains reachable over the correct BGP route, as learned by MetalLB FRR, after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that MetalLB FRR is configured for BGP routing and has recovered its routing state.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyMetallbEgressTrafficSegregation`) to validate the egress traffic segregation via MetalLB FRR BGP routes.
2. This typically involves deploying an application that uses MetalLB for egress traffic and verifying that its traffic is routed through the expected BGP routes after a graceful reboot.
3. The test then verifies network connectivity and performance between these workloads after a graceful reboot, confirming their resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable external connectivity via MetalLB FRR BGP routes after a graceful reboot.
