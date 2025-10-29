# Test Case Summary for 54572

Test case 54572 is located in tests/cnf/ran/powermanagement/tests/powersave.go and is named "Enables powersave at node level and then enable performance at node level".

## Goal

The goal of this test case is to verify that enabling powersave at the node level, and then enabling performance at the node level, correctly updates the kernel parameters.

## Test Setup

Prior to the test case, it skips if the cluster is not a Single Node OpenShift (SNO) cluster or if the node is not an "amd64" architecture. It ensures that a performance profile with a CPU set is available.

It does not require a git config set up.

## Test Steps

1. It patches the performance profile with workload hints to enable powersave and waits for the Machine Config Pool (MCP) to update.
2. It executes the `cat /proc/cmdline` command on the SNO node to get the kernel command line parameters.
3. It then verifies that the `intel_pstate=passive` kernel parameter is present in the output.
4. It also verifies that the `intel_pstate=disable` kernel parameter is NOT present in the output.


