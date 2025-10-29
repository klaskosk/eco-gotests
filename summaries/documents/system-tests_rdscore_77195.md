# Test Case Summary for 77195

Test case 77195 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies Multus-Tap CNI for rootless DPDK on the same node, single VF with multiple VLANs".

## Goal

The goal of this test case is to verify the functionality of Multus-Tap CNI for rootless DPDK on the same node, using a single Virtual Function (VF) with multiple VLANs.

## Test Setup

Prior to the test case, this test assumes a cluster with Multus-Tap CNI and rootless DPDK configured, with a single VF and multiple VLANs set up on the same node.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyRootlessDPDKOnTheSameNodeSingleVFMultipleVlans` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that rootless DPDK workloads can effectively use Multus-Tap CNI with a single VF and multiple VLANs for network connectivity on the same node.
