# Test Case Summary for CleanupUnexpectedAdmissionPodsCP

Test case CleanupUnexpectedAdmissionPodsCP is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Cleanup UnexpectedAdmission pods after KDump test on Control Plane node".

## Goal

The goal of this test case is to clean up any pods that are in an "UnexpectedAdmissionError" state after the KDump test on the Control Plane node.

## Test Setup

Prior to the test case, this test assumes that a KDump test has been executed on the Control Plane node, which might have left some pods in an "UnexpectedAdmissionError" state.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.CleanupUnexpectedAdmissionPodsCP` to perform the cleanup. The detailed steps are within this helper function, but the overall intent is to identify and delete pods that failed admission after the KDump test.
