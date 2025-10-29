# Test Case Summary for 47950

Test case 47950 is located in tests/cnf/ran/talm/tests/talm-precache.go and is named "tests for ocp cache with version".

## Goal

The goal of this test case is to verify OCP image precaching with a specified version.

## Test Setup

Prior to the test case, the following changes are needed:

- No specific setup changes are mentioned beyond the `AfterEach` cleanup which ensures resources on the hub are cleaned up.

It does not require a git config set up such that X.

## Test Steps

1. A CGU is created and a policy with a ClusterVersion CR (defining upgrade graph, channel, and version) is applied.
2. The `GetClusterVersionDefinition` helper function is used to retrieve the cluster version definition for "Version" from `Spoke1APIClient`.
3. The `SetupCguWithClusterVersion` helper function is used to set up the CGU with the retrieved cluster version.
4. Wait until the CGU pre-cache status for `Spoke1Name` is "Succeeded".
5. Verify that the new precache pod on `Spoke1APIClient` succeeded by checking its logs for "Image pre-cache done".
