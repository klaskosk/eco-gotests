# Test Case Summary for 64746

Test case 64746 is located in tests/cnf/ran/talm/tests/talm-precache.go and is named "tests custom image precaching using a PreCachingConfig CR".

## Goal

The goal of this test case is to verify custom image precaching using a `PreCachingConfig` Custom Resource (CR).

## Test Setup

Prior to the test case, the following changes are needed:

- The test first checks if the TALM version is 4.14 or newer. If not, the test is skipped.
- It identifies the PTP image used by `Spoke1APIClient` and stores it as `targetPrecacheImage`.
- Any existing instances of `targetPrecacheImage` are deleted from `Spoke1APIClient` master nodes.
- A `PreCachingConfig` CR is created on the hub with `SpaceRequired` set to "10 GiB", `ExcludePrecachePatterns` to an empty string array, and `AdditionalImages` including the `targetPrecacheImage`.
- The `AfterEach` hook ensures that any `PreCachingConfig` created on the hub is deleted and other test resources are cleaned up.

It does not require a git config set up such that X.

## Test Steps

1. A CGU is defined, incorporating the `PreCachingConfigRef` to the created `PreCachingConfig`.
2. The `GetClusterVersionDefinition` helper function is used to retrieve the cluster version definition for "Image" from `Spoke1APIClient`.
3. The `SetupCguWithClusterVersion` helper function is used to set up the CGU with the retrieved cluster version.
4. Wait until the CGU pre-cache status for `Spoke1Name` is "Succeeded".
5. Generate a list of precached images on `Spoke1APIClient`, specifically filtering for the `targetPrecacheImage`.
6. Check that the `targetPrecacheImage` is present in the list of precached images on at least one master node of `Spoke1APIClient`.
