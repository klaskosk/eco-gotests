# Test Case Summary for 45698

Test case 45698 is located in tests/system-tests/ran-du/tests/launch-workload-multiple-iter-loadavg.go and is named "Launch workload multiple times".

## Goal

The goal of this test case is to verify the stability and performance of the cluster by repeatedly launching and tearing down a workload, ensuring that all components (deployments, statefulsets, pods) become ready, PTP synchronization is maintained (if enabled), and node load average remains within acceptable limits.

## Test Setup

Prior to the test case, the following changes are needed:

- No specific changes are needed as the `BeforeAll` block handles workload preparation and cleanup for each iteration.

It does not require a git config set up.

## Test Steps

1. For `RanDuTestConfig.LaunchWorkloadIterations` iterations:
    a. If the test workload namespace exists, delete the existing workload using the shell method.
    b. If `RanDuTestConfig.TestWorkload.CreateMethod` is `TestWorkloadShellLaunchMethod`, launch the workload using the shell method.
    c. Wait for all deployment replicas to become ready in the test workload namespace.
    d. Wait for all statefulset replicas to become ready in the test workload namespace.
    e. Wait for all pods to become ready in the test workload namespace.
    f. If PTP is enabled, wait for 3 minutes and then check PTP status, asserting that it is in sync.
2. Observe node load average while the workload is running:
    a. For 30 iterations, with a 10-second sleep between each:
        i. Execute `cat /proc/loadavg` on the cluster.
        ii. Parse the first load average value from the command output.
        iii. Assert that the parsed load average is numerically less than `RanDuTestConfig.TestMultipleLaunchWorkloadLoadAvg`.
3. After all iterations, clean up the test workload resources by deleting the workload.
