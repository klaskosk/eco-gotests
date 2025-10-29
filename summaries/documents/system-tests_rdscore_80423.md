# Test Case Summary for 80423

Test case 80423 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies SR-IOV workloads on different nodes and same SR-IOV network post reboot".

## Goal

The goal of this test case is to verify that SR-IOV workloads, distributed across different nodes and utilizing the same SR-IOV network, remain functional and connected after an ungraceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that SR-IOV networks are configured across multiple nodes with workloads deployed using the same SR-IOV network.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifySRIOVWorkloadsOnDifferentNodes`) to validate the functionality and connectivity of SR-IOV workloads.
2. This typically involves performing network connectivity tests (e.g., ping, iperf) between SR-IOV-enabled pods on different nodes sharing the same SR-IOV network, demonstrating their resilience post-reboot.
3. The intent is to confirm the proper functioning of SR-IOV workloads after a disruptive event, especially in a multi-node, shared network environment.
