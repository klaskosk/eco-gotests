# Test Case Summary for 59948

Test case 59948 is located in tests/cnf/ran/talm/tests/talm-precache.go and is named "tests precache image filtering".

## Goal

The goal of this test case is to verify configurable filters for precache images, specifically excluding an image.

## Test Setup

Prior to the test case, the following changes are needed:

- The test first checks if the TALM version is 4.13 or newer. If not, the test is skipped.
- A ConfigMap named `PreCacheOverrideName` is created on the hub in `TestNamespace` with data `{"excludePrecachePatterns": "prometheus"}`. This configmap is designed to exclude images matching "prometheus" from precaching.
- The `AfterEach` hook ensures that any `PreCachingConfig` created on the hub is deleted and other test resources are cleaned up.

It does not require a git config set up such that X.

## Test Steps

1. A CGU is created and set up with an image filter.
2. The `GetClusterVersionDefinition` helper function is used to retrieve the cluster version definition for "Image" from `Spoke1APIClient`.
3. The `SetupCguWithClusterVersion` helper function is used to set up the CGU with the retrieved cluster version.
4. Wait until the CGU pre-cache status for `Spoke1Name` is "Succeeded".
5. Generate a list of precached images on `Spoke1APIClient`.
6. Check that the excluded image (prometheus) is *not* present in the list of precached images on at least one master node of `Spoke1APIClient`.
