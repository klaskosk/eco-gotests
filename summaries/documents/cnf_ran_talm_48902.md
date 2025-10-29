# Test Case Summary for 48902

Test case 48902 is located in tests/cnf/ran/talm/tests/talm-precache.go and is named "tests for precache operator with multiple sources".

## Goal

The goal of this test case is to verify image precaching for operators with multiple sources.

## Test Setup

Prior to the test case, the following changes are needed:

- The test verifies that `TalmPrecachePolicies` from the config are available on the hub. If not, it skips the test.
- It then copies policies with subscriptions and makes them non-compliant, appending a suffix to their names.

It does not require a git config set up such that X.

## Test Steps

1. A ClusterGroupUpgrade (CGU) is created with the operator upgrade policies.
2. Wait until the CGU pre-cache status for `Spoke1Name` is "Succeeded".
3. Verify that the image precache pod succeeded on the spoke by checking its logs for "Image pre-cache done".
