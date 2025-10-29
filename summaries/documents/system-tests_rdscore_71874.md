# Test Case Summary for 71874

Test case 71874 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies CephFS PVC is still accessible".

## Goal

The goal of this test case is to verify that a CephFS Persistent Volume Claim (PVC) remains accessible after a graceful cluster reboot, ensuring data persistence and application availability.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that a CephFS PVC was previously provisioned and is expected to be accessible.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyCephFSPVC`) to validate the accessibility of the CephFS PVC.
2. This typically involves deploying an application that attempts to read from or write to the CephFS PVC.
3. The test then verifies that the application can successfully access the CephFS PVC after a graceful reboot, confirming its resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable data access via CephFS PVC after a graceful reboot.
