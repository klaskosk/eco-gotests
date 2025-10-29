# Test Case Summary for 72826

Test case 72826 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Assert workload is reachable over IPv6 SPK ingress after TMM pod was deleted".

## Goal

The goal of this test case is to verify that a deployed workload remains reachable via its IPv6 address through the SPK Ingress after the SPK TMM (Traffic Management Microkernel) pods are deleted.

## Test Setup

Prior to the test case, it is assumed that:

- Existing SPK TMM pods (matching `tmmLabel`) are present in both `SPKConfig.SPKDataNS` and `SPKConfig.SPKDnsNS`.
- The SPK backend workload is set up as described in `SetupSPKBackendWorkload()` (ConfigMap, TCP service, and TCP deployment).

It does not require a git config set up.

## Test Steps

1. All SPK TMM data plane pods (matching `tmmLabel`) in `SPKConfig.SPKDataNS` are deleted.
2. All SPK TMM DNS plane pods (matching `tmmLabel`) in `SPKConfig.SPKDnsNS` are deleted.
3. The test waits for new TMM pods to become ready in both data and DNS namespaces.
4. The `AssertIPv6WorkloadURL` function is called to attempt accessing the configured IPv6 URL for the SPK-backed TCP workload and asserts that it returns a 200 OK status code, indicating successful reachability after the TMM pods were deleted and new ones came up.
