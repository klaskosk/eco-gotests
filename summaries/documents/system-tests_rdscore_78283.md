# Test Case Summary for 78283

Test case 78283 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify eIPv6 address assigned to the next available node after node reboot; fail-over".

## Goal

The goal of this test case is to verify the fail-over mechanism of eIPv6 addresses, ensuring that an eIPv6 address is reassigned to the next available node after a node reboot.

## Test Setup

Prior to the test case, this test assumes that EgressIP is configured in the cluster with an eIPv6 address, and that there are multiple nodes available for fail-over. The test will simulate a node reboot.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyEgressIPFailOverIPv6` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that EgressIP provides high availability for IPv6 addresses by reassigning them to healthy nodes after a failure.
