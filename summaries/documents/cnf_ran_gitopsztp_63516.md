# Test Case Summary for 63516

Test case 63516 is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-policies-app.go and is named "Reference non-existence source CR yaml file - verifies a proper error is returned in ArgoCD app when a non-existent source-cr is used in PGT".

## Goal

The goal of this test case is to verify that the Argo CD application returns a proper error when a PolicyGeneratorTemplate (PGT) references a non-existent source CR YAML file.

## Test Setup

Prior to the test case, the original policies app source is saved and then reset after the test. It requires ZTP version 4.14 or later. After each test, policies, service accounts, namespaces, and storage classes related to custom source CRs are deleted from the spoke cluster to ensure a clean state.

It does require a git config set up such that the policies app can be updated via GitOps.

## Test Steps

1. Check if the ZTP test path for a non-existent source CR file exists.
2. Update the Argo CD policies app with the Git path (`tsparams.ZtpTestPathCustomSourceNoCrFile`) and wait for synchronization.
3. Pull the Argo CD policies app.
4. Check the Argo CD conditions for the expected error message: "test/NoCustomCr.yaml is not found".
