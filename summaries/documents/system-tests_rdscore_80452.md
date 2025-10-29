# Test Case Summary for 80452

Test case 80452 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies SR-IOV workloads on the different nodes and different SR-IOV nets post graceful reboot".

## Goal

The goal of this test case is to verify that SR-IOV workloads, deployed on different nodes and utilizing different SR-IOV networks, remain functional and connected after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that SR-IOV networks are configured across different nodes with workloads deployed using different SR-IOV networks.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifySRIOVWorkloadsOnDifferentNodesDifferentNet`) to validate the functionality and connectivity of SR-IOV workloads.
2. This typically involves deploying SR-IOV enabled applications across different nodes, each configured to use different SR-IOV networks.
3. The test then verifies network connectivity and performance between these distributed SR-IOV workloads after a graceful reboot, confirming their resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable operation of SR-IOV workloads across different nodes with different SR-IOV networks after a graceful reboot.
