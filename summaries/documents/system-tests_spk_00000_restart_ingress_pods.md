# Test Case Summary for system-tests_spk_tests_no_id_restart_ingress_pods

Test case with no explicit ID, but labeled as "Restart SPK Ingress pods after hard reboot", is located in tests/system-tests/spk/tests/spk-suite.go.

## Goal

The goal of this test case is to restart SPK Ingress pods after a hard reboot, addressing potential issues with pods using cached data.

## Test Setup

This test case is part of the "Hard reboot" describe block and runs after a hard reboot has been performed. It assumes that SPK Ingress pods may be in a state requiring a restart.

It does not require a git config set up.

## Test Steps

1. All SPK Ingress data plane pods (matching `ingressDataLabel`) in `SPKConfig.SPKDataNS` are deleted.
2. All SPK Ingress DNS plane pods (matching `ingressDNSLabel`) in `SPKConfig.SPKDnsNS` are deleted.
3. The test waits for new Ingress pods to become ready in both data and DNS namespaces, ensuring they are functional after the restart.
