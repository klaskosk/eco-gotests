# Test Case Summary for SetupsEgressServiceWithLocalExternalTrafficPolicy

Test case SetupsEgressServiceWithLocalExternalTrafficPolicy is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Setups EgressService with Local ExternalTrafficPolicy".

## Goal

The goal of this test case is to set up EgressService with Local ExternalTrafficPolicy.

## Test Setup

Prior to the test case, this test assumes a cluster where EgressService and the necessary components for Local ExternalTrafficPolicy are available.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyEgressServiceWithLocalETP` to perform the setup. The detailed steps are within this helper function, but the overall intent is to configure EgressService to use `ExternalTrafficPolicy=Local` for the purpose of further testing, likely after an ungraceful reboot.
