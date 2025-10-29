# Test Case Summary for 79483

Test case 79483 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify EgressService with Local ExternalTrafficPolicy and sourceIPBy=Network".

## Goal

The goal of this test case is to verify EgressService with Local ExternalTrafficPolicy and `sourceIPBy=Network`.

## Test Setup

Prior to the test case, no specific changes are mentioned as required. This test likely assumes a pre-configured cluster with EgressService configured for `ExternalTrafficPolicy=Local` and `sourceIPBy=Network`.

It does not require a git config set up.

## Test Steps

1. The test calls `rdscorecommon.VerifyEgressServiceWithLocalETPSourceIPByNetwork` to perform the verification. The detailed steps are within this helper function, but the overall intent is to ensure EgressService functions correctly with `ExternalTrafficPolicy=Local` and `sourceIPBy=Network`.
