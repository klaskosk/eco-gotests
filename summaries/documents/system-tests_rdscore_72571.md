# Test Case Summary for 72571

Test case 72571 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies MACVLAN workloads on different nodes post graceful reboot".

## Goal

The goal of this test case is to verify that MACVLAN workloads, deployed across different nodes, remain functional and connected after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that MACVLAN networks are configured across multiple nodes with workloads deployed using these networks.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyMacVlanOnDifferentNodes`) to validate the functionality and connectivity of MACVLAN workloads.
2. This typically involves deploying MACVLAN-enabled applications across different nodes.
3. The test then verifies network connectivity and performance between these MACVLAN workloads after a graceful reboot, confirming their resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable operation of MACVLAN workloads across different nodes after a graceful reboot.
