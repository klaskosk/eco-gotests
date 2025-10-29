# Test Case Summary for 59946

Test case 59946 is located in tests/cnf/ran/talm/tests/talm-precache.go and is named "verifies CGU afterCompletion action executes on spoke2 when spoke1 is offline".

## Goal

The goal of this test case is to verify that the `afterCompletion` action of a ClusterGroupUpgrade (CGU) executes successfully on an online spoke (`Spoke2Name`), even when another spoke (`Spoke1Name`) in the batch is offline.

## Test Setup

Prior to the test case, the following changes are needed:

- The `BeforeAll` hook sets up a CGU with `RANConfig.Spoke1Name` (powered off) and `RANConfig.Spoke2Name` (available), with a batch size of 1. A timeout of 17 minutes is set for the remediation strategy.
- The CGU is updated to include an `afterCompletion` action that adds the label `talmcomplete` to clusters upon completion.
- The `AfterAll` hook cleans up resources on `Spoke2APIClient` and the hub, and deletes the `talmcomplete` label from `RANConfig.Spoke2Name`.

It requires a BMC configuration to be set up.

## Test Steps

1. Check `RANConfig.Spoke2Name` for the presence of the `talmcomplete` label using `helper.DoesClusterLabelExist`. Assert that the label is present.
2. Check `RANConfig.Spoke1Name` for the presence of the `talmcomplete` label using `helper.DoesClusterLabelExist`. Assert that the label is *not* present.
