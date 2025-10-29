# Test Case Summary for 72569

Test case 72569 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies MACVLAN workloads on the same node post hard reboot".

## Goal

The goal of this test case is to verify that MACVLAN workloads, deployed on the same node, remain functional and connected after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that MACVLAN networks are configured on the same node with workloads deployed using these networks.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyMacVlanOnSameNode`) to validate the functionality and connectivity of MACVLAN workloads.
2. This typically involves performing network connectivity tests (e.g., ping, iperf) between MACVLAN-enabled pods on the same node, demonstrating their resilience and proper network attachment post-reboot.
3. The intent is to confirm the proper functioning of MACVLAN workloads after a disruptive event, especially in a single-node environment.
