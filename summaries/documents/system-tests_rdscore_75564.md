# Test Case Summary for 75564

Test case 75564 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies IPVLAN workloads on the same node post hard reboot".

## Goal

The goal of this test case is to verify that IPVLAN workloads, deployed on the same node, remain functional and connected after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that IPVLAN networks are configured on the same node with workloads deployed using these networks.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyIPVlanOnSameNode`) to validate the functionality and connectivity of IPVLAN workloads.
2. This typically involves performing network connectivity tests (e.g., ping, iperf) between IPVLAN-enabled pods on the same node, demonstrating their resilience and proper network attachment post-reboot.
3. The intent is to confirm the proper functioning of IPVLAN workloads after a disruptive event, especially in a single-node environment.
