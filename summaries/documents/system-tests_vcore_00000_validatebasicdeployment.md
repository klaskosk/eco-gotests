# Test Case Summary for vCore Basic Deployment Suite

Test case `01_validate_basic_deployment.go` is located in `tests/system-tests/vcore/tests/01_validate_basic_deployment.go` and is part of the "vCore Basic Deployment Suite". It calls the functions `VerifyInitialDeploymentConfig()`, `VerifyCGroupDefault()`, and `VerifyPostDeploymentConfig()`.

## Goal

The goal of this test case is to verify the initial and post-deployment configurations of the vCore cluster, including cluster health, time synchronization, MachineConfigPool (MCP) deployments, node availability, and cgroup versioning.

## Test Setup

No specific manual changes are explicitly mentioned in the test file itself. The test suite is composed of three main verification functions:

- `VerifyInitialDeploymentConfig()`: This suite verifies the cluster's initial deployment status, including a healthy cluster status, time synchronization for master and worker nodes, and the availability of ODF, control-plane-worker, and user-plane-worker MCPs and nodes.
- `VerifyCGroupDefault()`: This suite verifies that cgroupv2 is the default cgroup mode for the cluster and tests the ability to switch between cgroupv1 and cgroupv2.
- `VerifyPostDeploymentConfig()`: This suite verifies post-deployment configurations such as Image Registry management state enablement, network policy configuration, SCC activation, SCTP module activation, and system reserved memory for master nodes.

It does not require a git config set up.

## Test Steps

**VerifyInitialDeploymentConfig() steps:**
1. Verify healthy cluster status by checking API URL availability, BareMetalHosts operational state, control-plane nodes count (expected 3), and that all master nodes are Ready.
2. Assert time sync was successfully applied for master nodes by verifying `/etc/chrony.conf` content.
3. Assert time sync was successfully applied for worker nodes by verifying `/etc/chrony.conf` content.
4. Verify ODF MCP was deployed and is in a ready state.
5. Verify full set of ODF nodes was deployed and are in a ready state.
6. Verify control-plane-worker MCP was deployed and is in a ready state.
7. Verify control-plane-worker nodes availability and ensure they are in a ready state.
8. Verify user-plane-worker MCP was deployed and is in a ready state.
9. Verify user-plane-worker nodes availability and ensure they are in a ready state.

**VerifyCGroupDefault() steps:**
1. Verify cgroupv2 is a default for the cluster deployment by checking the `nodes.config` 'cluster' object and the actual cgroup version configured for each node.
2. Verify that the cluster can be moved to cgroupv1 by changing the cluster cgroup mode to `CgroupModeV1`.

**VerifyPostDeploymentConfig() steps:**
1. Verify Image Registry management state is Enabled by changing `imageregistryconfig` management state to `Managed` and setting storage to `EmptyDir`.
2. Verify network policy configuration procedure by applying network policies from templates.
3. Verify SCC activation succeeded by applying `scc-config.yaml`.
4. Verify SCTP module activation succeeded by applying an `sctp-module.yaml` MachineConfig and verifying SCTP is active on control-plane nodes.
5. Verify system reserved memory for masters succeeded by applying `system-reserved-masters.yaml` MachineConfig.
