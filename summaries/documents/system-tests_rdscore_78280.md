# Test Case Summary for 78280

Test case 78280 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify eIPv4 address assigned to the next available node after node reboot; fail-over".

## Goal

The goal of this test case is to verify the fail-over mechanism of eIPv4 addresses, ensuring that an eIPv4 address is reassigned to the next available node after a node reboot.

## Test Setup

Prior to the test case, this test assumes that EgressIP is configured in the cluster with an eIPv4 address, and that there are multiple nodes available for fail-over. The test will simulate a node reboot.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyEgressIPFailOverIPv4` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure that EgressIP provides high availability for IPv4 addresses by reassigning them to healthy nodes after a failure.
