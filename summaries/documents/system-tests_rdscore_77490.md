# Test Case Summary for 77490

Test case 77490 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies Multus-Tap CNI for rootless DPDK pod workloads on the different nodes, multiple IP-VLANs".

## Goal

The goal of this test case is to verify the functionality of Multus-Tap CNI for rootless DPDK pod workloads when deployed on different nodes and utilizing multiple IP-VLANs.

## Test Setup

Prior to the test case, this test assumes a cluster with Multus-Tap CNI and rootless DPDK configured, with multiple IP-VLANs set up across different nodes.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyRootlessDPDKWorkloadsOnDifferentNodesMultipleIPVlans` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that rootless DPDK workloads can effectively use Multus-Tap CNI with multiple IP-VLANs for network connectivity across different nodes.
