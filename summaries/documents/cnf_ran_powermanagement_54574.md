# Test Case Summary for 54574

Test case 54574 is located in tests/cnf/ran/powermanagement/tests/powersave.go and is named "Enable powersave, and then enable high performance at node level, check power consumption with no workload pods."

## Goal

The goal of this test case is to verify the power settings of CPUs after enabling powersave and then high performance at the node level, specifically when no workload pods are running. It also checks the power consumption with no workload pods.

## Test Setup

Prior to the test case, it skips if the cluster is not a Single Node OpenShift (SNO) cluster or if the node is not an "amd64" architecture. It ensures that a performance profile with a CPU set is available. It defines test pod annotations for CPU load balancing, quota, IRQ load balancing, C-states, and CPU frequency governor.

It does not require a git config set up.

## Test Steps

1. Defines a QoS test pod with specified CPU and memory limits and the previously defined annotations. The `RuntimeClassName` is set to the performance profile name.
2. Creates and waits for the test pod to be running. It then verifies that the test pod has a QoS class of `Guaranteed`.
3. Executes a command within the test pod to get the `cpuset` of the process and parses the output to identify the CPUs used by the pod.
4. Verifies the power settings (resume latency and scaling governor) of the CPUs used by the pod, expecting "n/a" for `pm_qos_resume_latency_us` and "performance" for `scaling_governor`.
5. Determines the CPUs not assigned to the pod and verifies their power settings, expecting "0" for `pm_qos_resume_latency_us` and "performance" for `scaling_governor`.
6. Deletes the test pod and waits for its termination.
7. After the pod is deleted, it verifies that the CPUs previously assigned to the container revert to default powersave settings, expecting "0" for `pm_qos_resume_latency_us` and "performance" for `scaling_governor`.
8. There's also a `DeferCleanup` function to ensure the test pod is deleted even if the test fails.


