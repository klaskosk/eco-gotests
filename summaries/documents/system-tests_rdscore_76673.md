# Test Case Summary for 76673

Test case 76673 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify EgressService ingress with Local ExternalTrafficPolicy after graceful reboot".

## Goal

The goal of this test case is to verify that ingress traffic routing through EgressService, configured with Local ExternalTrafficPolicy, functions correctly after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and the cluster is in a state where EgressService with Local ExternalTrafficPolicy for ingress should be functional.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyEgressServiceETPLocalIngressConnectivity`) to validate the ingress traffic flow.
2. This typically involves deploying an application that receives ingress traffic and is exposed via an EgressService with `ExternalTrafficPolicy` set to `Local`.
3. The test then sends traffic to this application and verifies that it is correctly routed to the local node through the EgressService, confirming its resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the integrity of ingress traffic segregation and routing after a graceful reboot.
