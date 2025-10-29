# Test Case Summary for 72194

Test case 72194 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts workload reachable via IPv6 address".

## Goal

The goal of this test case is to verify that a deployed workload is reachable via its IPv6 address through the SPK Ingress after a hard reboot of the cluster.

## Test Setup

Prior to the test case, the following setup is performed:

- The SPK backend UDP workload is set up using `spkcommon.SetupSPKBackendUDPWorkload()`.
- The SPK backend TCP workload is set up using `spkcommon.SetupSPKBackendWorkload()`.
- A hard reboot of the cluster is performed by `spkcommon.VerifyHardRebootSuite()`, which cycles power to the nodes.
- Any stuck SPK pods are cleaned up by `spkcommon.CleanupStuckContainerPods()`.
- SPK Ingress pods are restarted by `spkcommon.RestartSPKIngressPods()`.

It does not require a git config set up.

## Test Steps

1. The test attempts to access the configured IPv6 URL for the SPK-backed TCP workload.
2. It asserts that the HTTP GET request to the IPv6 URL returns a 200 OK status code, indicating successful reachability of the workload via its IPv6 address after the hard reboot and subsequent cleanup/restarts.
