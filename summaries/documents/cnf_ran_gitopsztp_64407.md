# Test Case Summary for 64407

Test case 64407 is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-policies-app.go and is named "Verify source CR search path implementation - verifies custom and default source CRs can be included in the same policy".

## Goal

The goal of this test case is to verify that both custom and default source CRs can be successfully included and applied within the same PolicyGeneratorTemplate (PGT) policy, demonstrating proper search path implementation for source CRs.

## Test Setup

Prior to the test case, the original policies app source is saved and then reset after the test. It requires ZTP version 4.14 or later. Before each test, it verifies that a service account (`tsparams.CustomSourceCrName`) and a storage class (`tsparams.CustomSourceStorageClass`) do not exist on the spoke cluster. After each test, policies, service accounts, namespaces, and storage classes related to custom source CRs are deleted from the spoke cluster to ensure a clean state.

It does require a git config set up such that the policies app can be updated via GitOps.

## Test Steps

1. Check if the service account (`tsparams.CustomSourceCrName`) does not exist on the spoke.
2. Check if the storage class (`tsparams.CustomSourceStorageClass`) does not exist on the spoke.
3. Check if the ZTP test path for custom source search path exists.
4. Update the Argo CD policies app with the Git path (`tsparams.ZtpTestPathCustomSourceSearchPath`) and wait for synchronization.
5. Wait for the policy (`tsparams.CustomSourceCrPolicyName`) to exist on the hub cluster.
6. Wait for the policy to become `Compliant`.
7. Check that the service account exists on the spoke.
8. Check that the storage class exists on the spoke.
