# Test Case Summary for 81430

Test case 81430 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies rootless DPDK pod workloads on the different nodes, multiple IP-VLANs post hard reboot".

## Goal

The goal of this test case is to verify that rootless DPDK (Data Plane Development Kit) pod workloads, deployed on different nodes and utilizing multiple IP-VLANs, remain functional and performant after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that Multus-Tap CNI and rootless DPDK are configured with multiple IP-VLANs set up across different nodes.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyRootlessDPDKWorkloadsOnDifferentNodesMultipleIPVlans`) to validate the functionality of the rootless DPDK setup.
2. This typically involves deploying rootless DPDK applications within pods distributed across different nodes, each configured to use multiple IP-VLANs.
3. The test then verifies network connectivity and data plane performance between these distributed DPDK pods, confirming that DPDK is operational and resilient post-reboot in a multi-node, multi-IP-VLAN environment.
4. The intent is to ensure the integrity and performance of rootless DPDK workloads across different nodes and multiple IP-VLANs after a disruptive cluster event.
