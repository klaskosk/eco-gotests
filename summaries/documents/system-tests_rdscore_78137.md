# Test Case Summary for 78137

Test case 78137 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify EgressIP connectivity over IPv6 address after ungraceful reboot".

## Goal

The goal of this test case is to verify that EgressIP (eIP) connectivity using an IPv6 address is fully functional after an ungraceful cluster reboot, ensuring applications can egress the cluster using the assigned eIP.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and the cluster is in a state where EgressIP functionality for IPv6 addresses should be restored and operational.

It does not require a git config set up.

## Test Steps

1. The test utilizes a helper function (`rdscorecommon.VerifyEgressIPConnectivityIPv6`) to validate the EgressIP connectivity with an IPv6 address.
2. This typically involves deploying a workload, assigning an eIPv6, and then verifying that outbound traffic from this workload uses the designated eIPv6, demonstrating resilience post-reboot.
3. The intent is to confirm the proper functioning of IPv6 EgressIP after a disruptive event.
