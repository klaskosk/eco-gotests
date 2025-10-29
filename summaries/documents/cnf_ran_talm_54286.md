# Test Case Summary for 54286

Test case 54286 is located in tests/cnf/ran/talm/tests/talm-precache.go and is named "verifies precaching fails for one spoke and succeeds for the other".

## Goal

The goal of this test case is to verify that when one managed cluster is powered off and unavailable, precaching fails for that spoke but succeeds for the available spoke, unblocking the batch OCP upgrade.

## Test Setup

Prior to the test case, the following changes are needed:

- The `BeforeAll` hook verifies that `HubAPIClient`, `Spoke1APIClient`, and `Spoke2APIClient` are not nil, and that `BMCClient` is configured. If not, the test is skipped.
- `Spoke1APIClient` is powered off using `rancluster.PowerOffWithRetries`.
- The `AfterEach` hook ensures that resources on the hub are cleaned up.
- The `AfterAll` hook powers on `Spoke1APIClient` and waits for all its pods to be ready using `cluster.WaitForRecover`.

It requires a BMC configuration to be set up.

## Test Steps

1. A CGU is created and set up with `RANConfig.Spoke1Name` (the powered-off spoke) and `RANConfig.Spoke2Name` (the available spoke) as target clusters.
2. The `GetClusterVersionDefinition` helper function is used to retrieve the cluster version definition for "Both" from `Spoke2APIClient`.
3. The `SetupCguWithClusterVersion` helper function is used to set up the CGU with the retrieved cluster version.
4. Wait for the CGU's `CguPreCacheValidCondition` to be met, confirming that precache is valid.
5. Wait until the CGU pre-cache status for `RANConfig.Spoke2Name` is "Succeeded".
6. The CGU is enabled by setting `cguBuilder.Definition.Spec.Enable` to `true` and updating the CGU.
7. Wait until the CGU reports `CguPreCachePartialCondition`, indicating that one spoke failed precaching.
8. Check that the CGU reports `RANConfig.Spoke1Name` failed with "UnrecoverableError".
