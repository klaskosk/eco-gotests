# Test Case Summary for 79518

Test case 79518 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify EgressService with Cluster ExternalTrafficPolicy and sourceIPBy=Network after graceful reboot".

## Goal

The goal of this test case is to verify the functionality of EgressService with Cluster ExternalTrafficPolicy and `sourceIPBy=Network` after a graceful cluster reboot, ensuring egress traffic is correctly routed with network-based source IP.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and the cluster is in a state where EgressService with Cluster ExternalTrafficPolicy and `sourceIPBy=Network` should be functional.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyEgressServiceConnectivityETPClusterSourceIPByNetwork`) to validate the EgressService configuration.
2. This typically involves deploying an application that utilizes EgressService with `ExternalTrafficPolicy` set to `Cluster` and `sourceIPBy` set to `Network`, then verifying that its egress traffic is correctly routed through the cluster, maintaining connectivity with the specified source IP after the graceful reboot.
3. The intent is to confirm the resilience and proper operation of EgressService with Cluster ExternalTrafficPolicy and `sourceIPBy=Network` in a controlled cluster restart scenario.
