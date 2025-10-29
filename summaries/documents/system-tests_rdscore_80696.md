# Test Case Summary for 80696

Test case 80696 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies pod-level bonded workloads after bond interface recovering after both VFs failure".

## Goal

The goal of this test case is to verify the behavior and recovery of pod-level bonded workloads after both Virtual Functions (VFs) of a bonded interface experience failure and then recover.

## Test Setup

Prior to the test case, this test assumes a cluster with pod-level bonding configured with multiple VFs. The test will simulate a failure of both VFs and their subsequent recovery.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyPodLevelBondWorkloadsAfterBothVFsFailure` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that bonded workloads can gracefully handle and recover from a complete failure of all VFs in a bonded interface, maintaining network connectivity.
