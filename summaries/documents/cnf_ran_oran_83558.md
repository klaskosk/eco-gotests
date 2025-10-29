# Test Case Summary for 83558

Test case 83558 is located in tests/cnf/ran/oran/tests/oran-alarms.go and is named "retrieves a subscription from the API".

## Goal

The goal of this test case is to verify that an alarm subscription can be retrieved from the O2IMS API after it has been created.

## Test Setup

Prior to the test case, an `O2IMS API client` and an `Alertmanager API client` are created. The `spoke 1 cluster ID` is also retrieved. These are done in the `BeforeEach` block.

It does not require a git config set up such that X.

## Test Steps

1. A test subscription is created with a `ConsumerSubscriptionId` and a `Callback` URL.
2. The subscription is retrieved from the `O2IMS API` using its `AlarmSubscriptionId`.
3. It is verified that the retrieved subscription's `ConsumerSubscriptionId` matches the `ConsumerSubscriptionId` of the initially created subscription.
4. All subscriptions are listed and it is verified that the retrieved subscription is present in the list.
5. The test subscription is deleted.
