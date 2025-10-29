# Test Case Summary for 72785

Test case 72785 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts workload reachable via IPv4 UDP after hard reboot".

## Goal

The goal of this test case is to verify that a deployed UDP workload is reachable via its IPv4 address through the SPK Ingress after a hard reboot of the cluster.

## Test Setup

Prior to the test case, the following setup is performed:

- The SPK backend UDP workload is set up using `spkcommon.SetupSPKBackendUDPWorkload()` (UDP service and UDP deployment).
- A hard reboot of the cluster is performed by `spkcommon.VerifyHardRebootSuite()`, which cycles power to the nodes.
- Any stuck SPK pods are cleaned up by `spkcommon.CleanupStuckContainerPods()`.
- SPK Ingress pods are restarted by `spkcommon.RestartSPKIngressPods()`.

It does not require a git config set up.

## Test Steps

1. The test attempts to dial a UDP connection to the configured IPv4 UDP URL.
2. A UDP message with a unique timestamp is sent to the workload.
3. The test verifies that the sent UDP message is present in the logs of the UDP server pods, confirming successful reachability of the workload via its IPv4 UDP after the hard reboot and subsequent cleanup/restarts.
