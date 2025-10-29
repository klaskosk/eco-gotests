# Test Case Summary for 77488

Test case 77488 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies Multus-Tap CNI for rootless DPDK pod workloads on the different nodes, multiple MACVLANs".

## Goal

The goal of this test case is to verify the functionality of Multus-Tap CNI for rootless DPDK pod workloads when deployed on different nodes and utilizing multiple MACVLANs.

## Test Setup

Prior to the test case, this test assumes a cluster with Multus-Tap CNI and rootless DPDK configured, with multiple MACVLANs set up across different nodes.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyRootlessDPDKWorkloadsOnDifferentNodesMultipleMacVlans` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that rootless DPDK workloads can effectively use Multus-Tap CNI with multiple MACVLANs for network connectivity across different nodes.
