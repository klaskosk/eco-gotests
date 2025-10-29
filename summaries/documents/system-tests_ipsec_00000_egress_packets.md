# Test Case Summary for IpsecPacketsSnoEgress

Test case IpsecPacketsSnoEgress is located in `tests/system-tests/ipsec/tests/ipsec_packets_ocp_egress.go` and is named "Asserts packets originating from the SNO cluster are sent via IPSec".

## Goal

The goal of this test case is to verify that egress traffic originating from a Single Node OpenShift (SNO) cluster is correctly routed through the IPSec tunnel.

## Test Setup

Prior to the test case, the `BeforeSuite` in `ipsec_suite_test.go` sets up a test namespace. This test's `BeforeAll` hook performs the following:

1.  Retrieves node names and converts `NodePort` and `NodePortIncrement` from the `IpsecTestConfig` to integers.
2.  For each cluster node, it creates a Kubernetes Service and an associated `iperf3` workload (deployment) in the test namespace.
    *   The Service is created with a unique `NodePort` for each node.
    *   The `iperf3` workload is deployed on a specific node and configured to run the `iperf3ToolImage`.

This test does not require a specific git config set up.

## Test Steps

1.  For each node in the cluster:
    a.  An `iperf3` server is started asynchronously on a Security Gateway (SecGW) via SSH using `sshcommand.SSHCommand()`. The server listens on the configured `NodePort` and `SecGwServerIP`.
    b.  After a brief delay, the current IPSec tunnel packet counts (`OutBytes` and `InBytes`) are retrieved using `ipsectunnel.TunnelPackets()`.
    c.  An `iperf3` client is launched asynchronously within the `iperf3` workload pod on the cluster node, targeting the SecGW `iperf3` server using `iperf3workload.LaunchIperf3Command()`.
    d.  The test waits for both the `iperf3` client and server operations to complete and asserts that no errors occurred.
    e.  The IPSec tunnel packet counts are retrieved again (`packetsAfter`).
    f.  The test asserts that the number of `OutBytes` after the `iperf3` client transmission is greater than the `OutBytes` before, indicating that traffic was sent via the IPSec tunnel.
