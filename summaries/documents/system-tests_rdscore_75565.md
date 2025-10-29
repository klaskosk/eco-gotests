# Test Case Summary for 75565

Test case 75565 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies IPVLAN workloads on the same node post graceful reboot".

## Goal

The goal of this test case is to verify that IPVLAN workloads, deployed on the same node, remain functional and connected after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that IPVLAN networks are configured on the same node with workloads deployed using these networks.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyIPVlanOnSameNode`) to validate the functionality and connectivity of IPVLAN workloads.
2. This typically involves deploying IPVLAN-enabled applications on the same node.
3. The test then verifies network connectivity and performance between these IPVLAN workloads after a graceful reboot, confirming their resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable operation of IPVLAN workloads on the same node after a graceful reboot.
