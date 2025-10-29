# Test Case Summary for 73728

Test case 73728 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies NUMA-aware workload is available after graceful reboot".

## Goal

The goal of this test case is to verify that NUMA-aware workloads remain available and correctly scheduled after a graceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that a graceful cluster reboot has occurred, and that NUMA-aware workloads are deployed and configured to be resilient to such events.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyNROPWorkload`) to validate the state and availability of NUMA-aware workloads.
2. This typically involves deploying a NUMA-aware application and monitoring its pod status and resource allocation to ensure it is running on the expected NUMA nodes.
3. The test then verifies that the NUMA-aware application continues to function correctly after a graceful reboot, confirming its resilience and proper operation after a controlled cluster restart.
4. The intent is to ensure the continuous and reliable operation of NUMA-aware workloads after a graceful reboot.
