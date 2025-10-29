# Test Case Summary for 81424

Test case 81424 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies rootless DPDK on the same node, single VF with multiple VLANs post graceful reboot".

## Goal

The goal of this test case is to verify that rootless DPDK (Data Plane Development Kit) on the same node, utilizing a single Virtual Function (VF) with multiple VLANs, functions correctly after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that Multus-Tap CNI (Container Network Interface) and rootless DPDK are configured with a single VF and multiple VLANs on the same node.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyRootlessDPDKOnTheSameNodeSingleVFMultipleVlans`) to validate the functionality of the rootless DPDK setup.
2. This typically involves deploying a rootless DPDK application within a pod on the specified node, configured with a single VF and multiple VLANs.
3. The test then verifies network connectivity and data plane performance across the configured VLANs, confirming that DPDK is operational and resilient post-reboot.
4. The intent is to ensure the integrity and performance of rootless DPDK workloads in a multi-VLAN, single VF setup on the same node after a controlled cluster restart.
