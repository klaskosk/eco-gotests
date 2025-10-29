# Test Case Summary for 62260

Test case 62260 is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-policies-app.go and is named "Same source CR file name - verifies the custom source CR takes precedence over the default source CR with the same file name".

## Goal

The goal of this test case is to verify that when a custom source CR has the same file name as a default source CR, the custom CR takes precedence, leading to the creation of resources defined in the custom CR.

## Test Setup

Prior to the test case, the original policies app source is saved and then reset after the test. It requires ZTP version 4.14 or later. After each test, policies, service accounts, namespaces, and storage classes related to custom source CRs are deleted from the spoke cluster to ensure a clean state.

It does require a git config set up such that the policies app can be updated via GitOps.

## Test Steps

1. Check if the ZTP test path for replacing an existing custom source CR exists.
2. Update the Argo CD policies app with the Git path (`tsparams.ZtpTestPathCustomSourceReplaceExisting`) and wait for synchronization.
3. Wait for the policy (`tsparams.CustomSourceCrPolicyName`) to exist on the hub cluster.
4. Wait for the policy to become `Compliant`.
5. Check that the custom namespace (`tsparams.CustomSourceCrName`) exists, indicating the custom CR took precedence.
