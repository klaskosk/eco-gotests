# Test Case Summary for Setups EgressService with Local ExternalTrafficPolicy and sourceIPBy=Network (Graceful Reboot)

Test case "Setups EgressService with Local ExternalTrafficPolicy and sourceIPBy=Network" is located in tests/system-tests/rdscore/tests/00_validate_top_level.go within the "Graceful Cluster Reboot" context.

## Goal

The goal of this test case is to set up EgressService with Local ExternalTrafficPolicy and `sourceIPBy=Network` after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that a cluster where EgressService and the necessary components for Local ExternalTrafficPolicy with `sourceIPBy=Network` are available and have recovered their state.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyEgressServiceWithLocalETPSourceIPByNetwork`) to set up and validate the EgressService configuration.
2. This typically involves deploying an EgressService with `ExternalTrafficPolicy` set to `Local` and `sourceIPBy` set to `Network`, ensuring it can correctly route egress traffic with network-based source IP after a graceful reboot.
3. The intent is to confirm that the EgressService can be successfully configured and remains operational with `Local` external traffic policy and `sourceIPBy=Network` after a controlled cluster restart.
