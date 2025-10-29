# Test Case Summary for 76504

Test case 76504 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify EgressService with Local ExternalTrafficPolicy after ungraceful reboot".

## Goal

The goal of this test case is to verify the functionality of EgressService with Local ExternalTrafficPolicy after an ungraceful cluster reboot, ensuring egress traffic maintains source IP and is routed to local nodes.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and the cluster is in a state where EgressService with Local ExternalTrafficPolicy should be functional.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyEgressServiceWithLocalETP`) to validate the EgressService configuration with Local ExternalTrafficPolicy.
2. This typically involves verifying that traffic egresses the cluster from the same node it originated, preserving the source IP, demonstrating that the EgressService is operational post-reboot with local external traffic policy.
3. The intent is to confirm the resilience and proper functioning of egress traffic management with local traffic policy after a disruptive event.
