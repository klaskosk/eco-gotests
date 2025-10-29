# Test Case Summary for 73727

Test case 73727 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies NUMA-aware workload is available after ungraceful reboot".

## Goal

The goal of this test case is to verify that NUMA-aware workloads remain available and correctly scheduled after an ungraceful cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that NUMA-aware workloads are deployed and configured to be resilient to such events.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyNROPWorkload`) to validate the state and availability of NUMA-aware workloads.
2. This typically involves checking that the workloads are running on the expected NUMA nodes and are accessible, demonstrating the resilience of NUMA-aware scheduling and resource allocation post-reboot.
3. The intent is to confirm the proper functioning of NUMA-aware workloads after a disruptive event.
