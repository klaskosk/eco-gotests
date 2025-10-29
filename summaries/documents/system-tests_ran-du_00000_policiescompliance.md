# Test Case Summary for no_id_policiescompliance

Test case no_id_policiescompliance is located in tests/system-tests/ran-du/tests/du-ztp-policies-compliance.go and is named "Asserts all policies are compliant".

## Goal

The goal of this case is to verify that all ZTP policies are in a "Compliant" state within the cluster. This ensures that the cluster adheres to its defined Zero Touch Provisioning policies.

## Test Setup

Prior to the test case, the following changes are needed:

- Fetch the OpenShift Cluster name.

It does not require a git config set up.

## Test Steps

1. Fetch the cluster name.
2. List all policies across all namespaces, specifically targeting the cluster's namespace.
3. For each policy retrieved:
    a. Log the policy name.
    b. Assert that the `ComplianceState` of the policy is "Compliant".
    c. Log the compliance state of the policy.
