# Test Case Summary for 80424

Test case 80424 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies SR-IOV workloads on different nodes and same SR-IOV network post graceful reboot".

## Goal

The goal of this test case is to verify that SR-IOV workloads, distributed across different nodes and utilizing the same SR-IOV network, remain functional and connected after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that SR-IOV networks are configured across multiple nodes with workloads deployed using the same SR-IOV network.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifySRIOVWorkloadsOnDifferentNodes`) to validate the functionality and connectivity of SR-IOV workloads.
2. This typically involves deploying SR-IOV enabled applications across different nodes, each configured to use the same SR-IOV network.
3. The test then verifies network connectivity and performance between these distributed SR-IOV workloads after a graceful reboot, confirming their resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable operation of SR-IOV workloads across different nodes with the same SR-IOV network after a graceful reboot.
