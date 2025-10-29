# Test Case Summary for 80428

Test case 80428 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies SR-IOV workloads on the same node and same SR-IOV network post reboot".

## Goal

The goal of this test case is to verify that SR-IOV workloads, deployed on the same node and utilizing the same SR-IOV network, remain functional and connected after an ungraceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that SR-IOV networks are configured on the same node with workloads deployed using the same SR-IOV network.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifySRIOVWorkloadsOnSameNode`) to validate the functionality and connectivity of SR-IOV workloads.
2. This typically involves performing network connectivity tests (e.g., ping, iperf) between SR-IOV-enabled pods on the same node sharing the same SR-IOV network, demonstrating their resilience post-reboot.
3. The intent is to confirm the proper functioning of SR-IOV workloads after a disruptive event, especially in a single-node, shared network environment.
