# Test Case Summary for 71993

Test case 71993 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies CephRBD workload is deployable after soft reboot".

## Goal

The goal of this test case is to verify that a new CephRBD workload can be successfully deployed and utilize CephRBD storage after a graceful (soft) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and the cluster is in a state where CephRBD storage is operational for new deployments.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyCephRBDPVC`) to validate the deployment of a CephRBD workload.
2. This typically involves attempting to create a new pod with a CephRBD PVC and verifying that the pod starts successfully and can access the CephRBD storage.
3. The test then verifies that the newly deployed CephRBD workload functions correctly after a graceful reboot, confirming its resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable deployment of CephRBD workloads after a graceful reboot.
