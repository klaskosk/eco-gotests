# Test Case Summary for 71852

Test case 71852 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies CephFS workload is deployable after soft reboot".

## Goal

The goal of this test case is to verify that a new CephFS workload can be successfully deployed and utilize CephFS storage after a graceful (soft) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and the cluster is in a state where CephFS storage is operational for new deployments.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyCephFSPVC`) to validate the deployment of a CephFS workload.
2. This typically involves attempting to create a new pod with a CephFS PVC and verifying that the pod starts successfully and can access the CephFS storage.
3. The test then verifies that the newly deployed CephFS workload functions correctly after a graceful reboot, confirming its resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable deployment of CephFS workloads after a graceful reboot.
