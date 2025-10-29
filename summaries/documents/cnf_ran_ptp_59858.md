# Test Case Summary for 59858

Test case 59858 is located in tests/cnf/ran/ptp/tests/ptp-node-reboot.go and is named "should return to same stable status after ptp node soft reboot".

## Goal

The goal of this test case is to verify that the PTP system returns to a stable (locked) status after a soft reboot of a PTP node.

## Test Setup

Prior to the test case, the following steps are performed:

1.  Determine if the cluster is SNO using `rancluster.IsSNO`.
2.  Select a node to reboot by listing PTP daemon set pods using `pod.List` and extracting the node name from the first pod.
3.  Record the current time (`rebootTime`).
4.  Soft reboot the selected node using `cluster.ExecCmdWithStdoutWithRetries` with the command "sudo systemctl reboot".
5.  Wait for the node to recover:
    a.  If SNO, wait for cluster recovery using `cluster.WaitForRecover`.
    b.  If not SNO, pull the node object using `nodes.Pull`, wait for it to become `NotReady` using `rebootedNode.WaitUntilNotReady`, then wait for it to become `Ready` using `rebootedNode.WaitUntilReady`, and finally wait for all pods on the rebooted node to be healthy using `pod.WaitForPodsInNamespacesHealthy`.

It does not require a git config set up.

## Test Steps

1.  Create a Prometheus API client using `querier.CreatePrometheusAPIForCluster`.
2.  Assert that all clocks on the rebooted node are in a `LOCKED` state using `metrics.AssertQuery`, with a stable duration of 10 seconds and a timeout of 10 minutes.
