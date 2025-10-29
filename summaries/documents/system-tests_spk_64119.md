# Test Case Summary for 64119

Test case 64119 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts workload reachable via IPv4 address".

## Goal

The goal of this test case is to verify that a deployed workload is reachable via its IPv4 address through the SPK Ingress.

## Test Setup

Prior to the test case, the following setup is performed:

- A ConfigMap named "spk-cm" is created in the SPK namespace with default web page content.
- A TCP service named "f5-hello-world" is created in the SPK namespace.
- A TCP deployment named "spk-hello-world" is created with one replica.

It does not require a git config set up.

## Test Steps

1. The test attempts to access the configured IPv4 URL for the SPK-backed TCP workload.
2. It asserts that the HTTP GET request to the IPv4 URL returns a 200 OK status code, indicating successful reachability of the workload via its IPv4 address.
