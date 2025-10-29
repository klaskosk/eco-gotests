# Test Case Summary for 54571

Test case 54571 is located in tests/cnf/ran/powermanagement/tests/powersave.go and is named "Verifies expected kernel parameters with no workload hints specified in PerformanceProfile".

## Goal

The goal of this test case is to verify that the expected kernel parameters are present when no WorkloadHints are specified in the PerformanceProfile.

## Test Setup

Prior to the test case, it skips if the cluster is not a Single Node OpenShift (SNO) cluster, if the node is not an "amd64" architecture, or if WorkloadHints are already present in the PerformanceProfile. It ensures that the PerformanceProfile does not include WorkloadHints.

It does not require a git config set up.

## Test Steps

1. It retrieves the current PerformanceProfile and checks if WorkloadHints are already present. If they are, the test is skipped.
2. It executes the `cat /proc/cmdline` command on the SNO node to get the kernel command line parameters.
3. It then verifies that the following required kernel parameters are present in the output:
    - `nohz_full=[0-9,-]+`
    - `tsc=nowatchdog`
    - `nosoftlockup`
    - `nmi_watchdog=0`
    - `mce=off`
    - `skew_tick=1`
    - `intel_pstate=disable`


