# Test Case Summary for 79515

Test case 79515 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify EgressService with Local ExternalTrafficPolicy and sourceIPBy=Network after ungraceful reboot".

## Goal

The goal of this case is to verify the EgressService with Local ExternalTrafficPolicy and `sourceIPBy=Network` after an ungraceful cluster reboot, ensuring proper egress with network-based source IP and local routing.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and the cluster is in a state where EgressService with Local ExternalTrafficPolicy and `sourceIPBy=Network` should be functional.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyEgressServiceWithLocalETPSourceIPByNetwork`) to validate the EgressService configuration.
2. This typically involves verifying that egress traffic is routed locally, preserving source IP, and using the network interface for source IP identification post-reboot.
3. The intent is to confirm the resilience and proper functioning of egress traffic management under specific local and network-based source IP policies after a disruptive event.
