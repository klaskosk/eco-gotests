# Test Case Summary for 54355

Test case 54355 is located in tests/cnf/ran/gitopsztp/tests/ztp-generator.go and is named "Generation of CRs for a single site from ztp container - generates and installs time crs, manifests, and policies, and verifies they are present".

## Goal

The goal of this test case is to verify that the ZTP generator container can successfully generate install-time Custom Resources (CRs), manifests, and policies for a single site, and that these generated files are present in the expected output directories.

## Test Setup

Prior to the test case, the current user is retrieved to determine the site config path (typically `/home/<user>/site-configs`). The test asserts that this path exists. After each test, the generated manifests and policies in `siteconfig/out` and `policygentemplates/out` are deleted to ensure a clean state.

It does not require a git config set up.

## Test Steps

1. Generate the install-time CRs and manifests using `podman run` with the `RANConfig.ZtpSiteGenerateImage` and `generator install -E /resources/` command.
2. Validate that the CRs and manifests were created by checking if at least 9 files exist in each site directory under `siteconfig/out/generated_installCRs/`.
3. Generate the policies using `podman run` with the `RANConfig.ZtpSiteGenerateImage` and `generator config .` command.
4. Validate that the policies were created by checking if at least 3 subdirectories (common, group DU, site) exist under `policygentemplates/out/generated_configCRs/`.
5. For each generated policy file, unmarshal its content as YAML and verify that its `kind` field is one of the expected values: "Policy", "PlacementRule", or "PlacementBinding".
