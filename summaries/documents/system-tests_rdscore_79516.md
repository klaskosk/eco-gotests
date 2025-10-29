# Test Case Summary for 79516

Test case 79516 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify EgressService ingress with Local ExternalTrafficPolicy and sourceIPBy=Network after ungraceful reboot".

## Goal

The goal of this test case is to verify that ingress traffic routing through EgressService, configured with Local ExternalTrafficPolicy and `sourceIPBy=Network`, functions correctly after an ungraceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and the cluster is in a state where EgressService with Local ExternalTrafficPolicy and `sourceIPBy=Network` for ingress should be functional.

It does not require a git config set up.

## Test Steps

1. The test utilizes a helper function (`rdscorecommon.VerifyEgressServiceIngressWithLocalETPSourceIPByNetwork`) to validate the ingress traffic flow through EgressService.
2. This typically involves sending traffic to a service and verifying that it is routed to the correct pods on local nodes, with source IP identified by the network, demonstrating the resilience of ingress functionality.
3. The intent is to confirm the proper functioning of ingress traffic management with local external traffic policy and network-based source IP after a disruptive event.
