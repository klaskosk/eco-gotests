# Test Case Summary for 78135

Test case 78135 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify eIPv6 address from the list of defined used for the assigned pods in a single eIP namespace".

## Goal

The goal of this test case is to verify that an eIPv6 address from a defined list is used for assigned pods in a single EgressIP (eIP) namespace, with balanced eIP traffic across three nodes.

## Test Setup

Prior to the test case, this test assumes that EgressIP is configured in the cluster with an eIPv6 address list, and that pods are assigned in a single eIP namespace across three nodes.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyEgressIPOneNamespaceThreeNodesBalancedEIPTrafficIPv6` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that eIPv6 addresses are correctly assigned and utilized by pods in a single eIP namespace, with traffic balanced across specified nodes.
