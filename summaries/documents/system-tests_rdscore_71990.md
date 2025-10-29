# Test Case Summary for 71990

Test case 71990 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies CephRBD PVC is still accessible".

## Goal

The goal of this test case is to verify that a CephRBD Persistent Volume Claim (PVC) remains accessible after an ungraceful cluster reboot, ensuring data persistence and application availability.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that a CephRBD PVC was previously provisioned and is expected to be accessible.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyCephRBDPVC`) to validate the accessibility of the CephRBD PVC.
2. This typically involves mounting the PVC to a pod and performing read/write operations to ensure data integrity and connectivity, demonstrating the resilience of CephRBD storage post-reboot.
3. The intent is to confirm the proper functioning of CephRBD storage after a disruptive event.
