# Test Case Summary for 71992

Test case 71992 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies CephRBD workload is deployable after hard reboot".

## Goal

The goal of this test case is to verify that a new CephRBD workload can be successfully deployed and utilize CephRBD storage after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and the cluster is in a state where CephRBD storage is operational for new deployments.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyCephRBDPVC`) to validate the deployment of a CephRBD workload.
2. This typically involves creating a new pod that mounts and accesses a CephRBD PVC, ensuring that the storage system is fully functional and new workloads can consume it post-reboot.
3. The intent is to confirm the proper functioning and recoverability of CephRBD for new deployments after a disruptive event.
