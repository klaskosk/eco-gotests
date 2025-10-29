# Test Case Summary for 83557

Test case 83557 is located in tests/cnf/ran/oran/tests/oran-alarms.go and is named "filters alarm subscriptions from the API".

## Goal

The goal of this test case is to verify that alarm subscriptions can be filtered from the O2IMS API based on specific criteria.

## Test Setup

Prior to the test case, an `O2IMS API client` and an `Alertmanager API client` are created. The `spoke 1 cluster ID` is also retrieved. These are done in the `BeforeEach` block.

It does not require a git config set up such that X.

## Test Steps

1. Two test subscriptions are created with unique `ConsumerSubscriptionId` and `Callback` URLs.
2. Subscriptions are filtered from the `O2IMS API` using a filter based on the first `ConsumerSubscriptionId`.
3. It is verified that the filtered results contain the first subscription but not the second subscription.
4. Both test subscriptions are deleted.
