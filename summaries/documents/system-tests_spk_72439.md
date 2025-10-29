# Test Case Summary for 72439

Test case 72439 is located in tests/system-tests/spk/tests/spk-suite.go and is named "Asserts workload reachable via IPv6 address after application recreation".

## Goal

The goal of this test case is to verify that a deployed workload remains reachable via its IPv6 address through the SPK Ingress after the application has been re-created.

## Test Setup

Prior to the test case, the following setup is performed:

- A ConfigMap named "spk-cm" is created in the SPK namespace with default web page content.
- A TCP service named "f5-hello-world" is created in the SPK namespace.
- A TCP deployment named "spk-hello-world" is created with one replica.

It does not require a git config set up.

## Test Steps

1. The `SetupSPKBackendWorkload()` function is called, which re-creates the backend workload (ConfigMap, Service, and Deployment).
2. The test attempts to access the configured IPv6 URL for the SPK-backed TCP workload.
3. It asserts that the HTTP GET request to the IPv6 URL returns a 200 OK status code, indicating successful reachability of the workload via its IPv6 address after recreation.
