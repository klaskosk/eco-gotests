# Test Case Summary for 48903

Test case 48903 is located in tests/cnf/ran/talm/tests/talm-precache.go and is named "tests for ocp cache with image".

## Goal

The goal of this test case is to verify OCP image precaching with an explicit image URL.

## Test Setup

Prior to the test case, the following changes are needed:

- The test finds a Prometheus pod image on `Spoke1APIClient` to use as the `excludedPreCacheImage`.
- It then attempts to delete any existing instances of this image from `Spoke1APIClient` master nodes to ensure a clean state.
- The `AfterEach` hook ensures that any `PreCachingConfig` created on the hub is deleted and other test resources are cleaned up.

It does not require a git config set up such that X.

## Test Steps

1. A CGU is created and a policy with a ClusterVersion CR (defining the upgrade graph, channel, and version) is applied.
2. The `GetClusterVersionDefinition` helper function is used to retrieve the cluster version definition for "Image" from `Spoke1APIClient`.
3. The `SetupCguWithClusterVersion` helper function is used to set up the CGU with the retrieved cluster version.
4. Wait until the CGU pre-cache status for `Spoke1Name` is "Succeeded".
5. Verify that the new precache pod on `Spoke1APIClient` succeeded by checking its logs for "Image pre-cache done".
6. Generate a list of precached images on `Spoke1APIClient`.
7. Check that the `excludedPreCacheImage` (the Prometheus pod image found in the `BeforeEach`) is present in the list of precached images on at least one master node of `Spoke1APIClient`.
