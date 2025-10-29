# Test Case Summary for IpsecTunnelAtDeployment

Test case IpsecTunnelAtDeployment is located in `tests/system-tests/ipsec/tests/ipsec_tunnel_deployment.go` and is named "Asserts the IPSec tunnel connected successfully at OCP deployment".

## Goal

The goal of this test case is to verify that the IPSec tunnel is successfully connected on all cluster nodes after OCP deployment.

## Test Setup

Prior to the test case, the common `BeforeSuite` in `ipsec_suite_test.go` ensures a dedicated test namespace (`ipsec-test`) is created. This namespace is labeled as privileged. No specific git config is required.

## Test Steps

1.  Retrieve the names of all nodes in the cluster using `ipsecinittools.GetNodeNames()`.
2.  For each node, execute the `ipsec show` command via a debug pod using `ipsectunnel.TunnelConnected()`.
3.  Verify that the `ipsec show` command returns output, indicating that the IPSec tunnel is connected on that node. If the output is empty, the test fails, signifying the tunnel is not established.
