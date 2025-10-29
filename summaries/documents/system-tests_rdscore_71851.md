# Test Case Summary for 71851

Test case 71851 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies CephFS workload is deployable after hard reboot".

## Goal

The goal of this test case is to verify that a new CephFS workload can be successfully deployed and utilize CephFS storage after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and the cluster is in a state where CephFS storage is operational for new deployments.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyCephFSPVC`) to validate the deployment of a CephFS workload.
2. This typically involves creating a new pod that mounts and accesses a CephFS PVC, ensuring that the storage system is fully functional and new workloads can consume it post-reboot.
3. The intent is to confirm the proper functioning and recoverability of CephFS for new deployments after a disruptive event.
