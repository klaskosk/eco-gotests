# Test Case Summary for 72778

Test case 72778 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts workload reachable via IPv6 UDP".

## Goal

The goal of this test case is to verify that a deployed UDP workload is reachable via its IPv6 address through the SPK Ingress.

## Test Setup

Prior to the test case, the following setup is performed:

- A UDP service named "f5-udp-svc" is created in the SPK namespace with NodePort 31225 and dual-stack IP families.
- A UDP deployment named "udp-mock-server" is created with one replica, running a UDP server on port 8080.

It does not require a git config set up.

## Test Steps

1. The test dials a UDP connection to the configured IPv6 UDP URL.
2. A UDP message with a unique timestamp is sent to the workload.
3. The test verifies that the sent UDP message is present in the logs of the UDP server pods, confirming successful reachability via IPv6 UDP.
