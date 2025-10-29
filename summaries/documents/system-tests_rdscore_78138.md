# Test Case Summary for 78138

Test case 78138 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify EgressIP connectivity over IPv6 address after graceful reboot".

## Goal

The goal of this test case is to verify that EgressIP (eIP) connectivity using an IPv6 address is fully functional after a graceful cluster reboot, ensuring applications can egress the cluster using the assigned eIP.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and the cluster is in a state where EgressIP functionality for IPv6 addresses should be restored and operational.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyEgressIPConnectivityIPv6`) to validate the EgressIP connectivity.
2. This typically involves deploying an application that utilizes an EgressIP with an IPv6 address.
3. The test then verifies that this application can successfully egress the cluster using the assigned IPv6 EgressIP, confirming its resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable external connectivity of applications via IPv6 EgressIP after a graceful reboot.
