# Test Case Summary for 53681

Test case 53681 is located in tests/cnf/ran/containernshide/tests/containernshide.go and is named "should not have kubelet and crio using the same inode as systemd".

## Goal

The goal of this test case is to verify that kubelet and crio processes are not using the same mount namespace as systemd, ensuring that the container namespaces are properly hidden on cluster nodes.

## Test Setup

This test case runs on a cluster with at least one spoke node (Spoke1APIClient is used). It relies on the `cluster.ExecCmdWithStdoutWithRetries` function to execute commands on the cluster nodes.

It does not require a git config set up.

## Test Steps

1. Get the mount namespace inode of the `systemd` process on all cluster nodes using `readlink /proc/1/ns/mnt`.
2. Get the mount namespace inode of the `kubelet` process on all cluster nodes using `readlink /proc/$(pidof kubelet)/ns/mnt`.
3. Get the mount namespace inode of the `crio` process on all cluster nodes using `readlink /proc/$(pidof crio)/ns/mnt`.
4. Verify that the number of collected inodes for `systemd`, `kubelet`, and `crio` are the same across all nodes.
5. For each host, verify that the `kubelet` inode matches the `crio` inode. This ensures that `kubelet` and `crio` are in the same mount namespace (as expected for container processes).
6. For each host, verify that the `systemd` inode does NOT match the `kubelet` inode, and the `systemd` inode does NOT match the `crio` inode. This confirms that `kubelet` and `crio` are running in a different mount namespace than `systemd`, indicating successful container namespace hiding.
