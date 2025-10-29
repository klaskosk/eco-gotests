# Test Case Summary for 'steadyworkload' scenario

Test case "Checks power usage for 'steadyworkload' scenario" is located in tests/cnf/ran/powermanagement/tests/powersave.go and is named "Checks power usage for 'steadyworkload' scenario".

## Goal

The goal of this test case is to collect and persist power usage metrics for a 'steadyworkload' scenario.

## Test Setup

Prior to the test case, it skips if BMC (Baseboard Management Controller) configuration is not set. It parses the `RANConfig.MetricSamplingInterval` to determine the sampling interval for metrics and retrieves the power state to be used as a tag for the metric from the performance profile.

It does not require a git config set up.

## Test Steps

1. It parses the `RANConfig.WorkloadDuration` to determine the duration for the 'steadyworkload' scenario.
2. It calls `collect.CollectPowerMetricsWithSteadyWorkload` to collect power metrics for the specified duration, sampling interval, power state, performance profile, and node name.
3. It then iterates through the collected metrics (component map) and prints each metric name and value to the Ginkgo report for further processing in the pipeline.


