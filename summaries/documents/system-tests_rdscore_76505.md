# Test Case Summary for 76505

Test case 76505 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify EgressService with Cluster ExternalTrafficPolicy after graceful reboot".

## Goal

The goal of this test case is to verify the functionality of EgressService with Cluster ExternalTrafficPolicy after a graceful cluster reboot, ensuring egress traffic is correctly routed.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and the cluster is in a state where EgressService with Cluster ExternalTrafficPolicy should be functional.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyEgressServiceConnectivityETPCluster`) to validate the EgressService configuration.
2. This typically involves deploying an application that utilizes EgressService with `ExternalTrafficPolicy` set to `Cluster` and then verifying that its egress traffic is correctly routed through the cluster, maintaining connectivity after the graceful reboot.
3. The intent is to confirm the resilience and proper operation of EgressService with Cluster ExternalTrafficPolicy in a controlled cluster restart scenario.
