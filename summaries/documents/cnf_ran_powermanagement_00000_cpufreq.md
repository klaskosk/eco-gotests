# Test Case Summary for CPU Frequency Tuning

Test case `cpufreq.go` is located in `tests/cnf/ran/powermanagement/tests/cpufreq.go` and is named "CPU frequency tuning tests change the core frequencies of isolated and reserved cores".

## Goal

The goal of this test case is to verify that the CPU frequencies of isolated and reserved cores are correctly set on the Device Under Test (DUT) when configured via a PerformanceProfile.

## Test Setup

Prior to the test case, the following changes are needed:

- The test uses `helper.GetPerformanceProfileWithCPUSet()` to retrieve an existing performance profile that has both reserved and isolated CPU sets defined. This function iterates through all performance profiles until one matching the criteria is found.
- It then parses the isolated and reserved CPU sets from the retrieved performance profile to identify an isolated core ID and a reserved core ID.
- It records the original CPU frequencies for both the isolated and reserved cores by calling `getCPUFreq` on each core, to be restored in the `AfterEach` block.

It does not require a git config set up.

## Test Steps

1. The test first checks if the OpenShift Container Platform (OCP) version is 4.16 or higher, skipping the test if it's not.
2. It then calls `helper.SetCPUFreq` with the performance profile and desired isolated and reserved core frequencies. This function updates the `PerformanceProfile` with the new `HardwareTuning` specifications. It then verifies that the CPU frequencies have been updated on the spoke cluster by executing `cat /sys/devices/system/cpu/cpufreq/policy<coreID>/scaling_max_freq` commands on the isolated and reserved CPUs and comparing the output to the desired frequencies, retrying until the frequencies match or a timeout occurs.
3. In the `AfterEach` block, the CPU frequencies are reverted to their original settings using `helper.SetCPUFreq` with the initially recorded original frequencies.
