# Test Case Summary for 30020

Test case 30020 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies ungraceful cluster reboot".

## Goal

The goal of this test case is to verify that the cluster can recover and function correctly after an ungraceful reboot.

## Test Setup

Prior to the test case, this test assumes a running cluster where an ungraceful reboot can be simulated.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyUngracefulReboot` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that critical cluster components and deployed workloads are resilient to sudden power loss or system crashes and can return to a healthy state.
