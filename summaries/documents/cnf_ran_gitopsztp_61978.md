# Test Case Summary for 61978

Test case 61978 is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-policies-app.go and is named "Create a new source CR in the user GIT repository - verifies new CR kind that does not exist in ztp container image can be created via custom source-cr".

## Goal

The goal of this test case is to verify that a new custom resource (CR) kind, not inherently present in the ZTP container image, can be successfully created and managed through a custom source CR within the user's Git repository via Argo CD GitOps.

## Test Setup

Prior to the test case, the original policies app source is saved and then reset after the test. It requires ZTP version 4.10 or later. After each test, policies, service accounts, namespaces, and storage classes related to custom source CRs are deleted from the spoke cluster to ensure a clean state.

It does require a git config set up such that the policies app can be updated via GitOps.

## Test Steps

1. Verify that the service account (`tsparams.CustomSourceCrName`) does not exist on the spoke cluster before the test begins.
2. Check if the ZTP test path for a new custom source CR exists.
3. Update the Argo CD policies app with the new custom source CR Git path (`tsparams.ZtpTestPathCustomSourceNewCr`) and wait for synchronization.
4. Wait for the policy (`tsparams.CustomSourceCrPolicyName`) to exist on the hub cluster.
5. Wait for the policy to become `Compliant`.
6. Wait for the service account (`tsparams.CustomSourceCrName`) to exist on the spoke cluster.
