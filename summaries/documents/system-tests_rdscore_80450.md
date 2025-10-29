# Test Case Summary for 80450

Test case 80450 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies SR-IOV workloads on same node and different SR-IOV nets post reboot".

## Goal

The goal of this test case is to verify that SR-IOV workloads, deployed on the same node and utilizing different SR-IOV networks, remain functional and connected after an ungraceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that SR-IOV networks are configured on the same node with workloads deployed using different SR-IOV networks.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifySRIOVWorkloadsOnSameNodeDifferentNet`) to validate the functionality and connectivity of SR-IOV workloads.
2. This typically involves performing network connectivity tests (e.g., ping, iperf) between SR-IOV-enabled pods on the same node using different SR-IOV networks, demonstrating their resilience post-reboot.
3. The intent is to confirm the proper functioning of SR-IOV workloads after a disruptive event, especially in a single-node, multi-network environment.
