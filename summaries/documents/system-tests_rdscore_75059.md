# Test Case Summary for 75059

Test case 75059 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies IPVLAN workloads on different nodes post graceful reboot".

## Goal

The goal of this test case is to verify that IPVLAN workloads, deployed across different nodes, remain functional and connected after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that IPVLAN networks are configured across multiple nodes with workloads deployed using these networks.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyIPVlanOnDifferentNodes`) to validate the functionality and connectivity of IPVLAN workloads.
2. This typically involves deploying IPVLAN-enabled applications across different nodes.
3. The test then verifies network connectivity and performance between these IPVLAN workloads after a graceful reboot, confirming their resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable operation of IPVLAN workloads across different nodes after a graceful reboot.
