# Test Case Summary for 72827

Test case 72827 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Assert workload is reachable over IPv4 SPK UDP ingress after TMM pod is deleted".

## Goal

The goal of this test case is to verify that a deployed UDP workload remains reachable via its IPv4 address through the SPK Ingress after the SPK TMM (Traffic Management Microkernel) pods are deleted.

## Test Setup

Prior to the test case, it is assumed that:

- Existing SPK TMM pods (matching `tmmLabel`) are present in both `SPKConfig.SPKDataNS` and `SPKConfig.SPKDnsNS`.
- The SPK backend UDP workload is set up as described in `SetupSPKBackendUDPWorkload()` (UDP service and UDP deployment).

It does not require a git config set up.

## Test Steps

1. All SPK TMM data plane pods (matching `tmmLabel`) in `SPKConfig.SPKDataNS` are deleted.
2. All SPK TMM DNS plane pods (matching `tmmLabel`) in `SPKConfig.SPKDnsNS` are deleted.
3. The test waits for new TMM pods to become ready in both data and DNS namespaces.
4. The `VerifySPKIngressUDPviaIPv4` function is called to attempt sending a UDP message to the configured IPv4 UDP URL and verifies its presence in the UDP server logs, indicating successful reachability after the TMM pods were deleted and new ones came up.
