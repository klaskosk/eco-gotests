# Test Case Summary for 79519

Test case 79519 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify EgressService with Local ExternalTrafficPolicy and sourceIPBy=Network after graceful reboot".

## Goal

The goal of this test case is to verify the EgressService with Local ExternalTrafficPolicy and `sourceIPBy=Network` after a graceful cluster reboot, ensuring proper egress with network-based source IP and local routing.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and the cluster is in a state where EgressService with Local ExternalTrafficPolicy and `sourceIPBy=Network` should be functional.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyEgressServiceConnectivityETPLocalSourceIPByNetwork`) to validate the EgressService configuration.
2. This typically involves deploying an application that utilizes EgressService with `ExternalTrafficPolicy` set to `Local` and `sourceIPBy` set to `Network`, then verifying that its egress traffic is correctly routed through the local node, maintaining connectivity with the specified source IP after the graceful reboot.
3. The intent is to confirm the resilience and proper operation of EgressService with Local ExternalTrafficPolicy and `sourceIPBy=Network` in a controlled cluster restart scenario.
