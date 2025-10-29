# Test Case Summary for 83561

Test case 83561 is located in tests/cnf/ran/oran/tests/oran-alarms.go and is named "ensures reliability of the alarms service".

## Goal

The goal of this test case is to ensure the reliability of the alarm service by concurrently sending multiple alerts and verifying that all corresponding notifications are received by a subscriber.

## Test Setup

Prior to the test case, an `O2IMS API client` and an `Alertmanager API client` are created. The `spoke 1 cluster ID` is also retrieved. These are done in the `BeforeEach` block.

It does not require a git config set up such that X.

## Test Steps

1. A test subscription is created with a `ConsumerSubscriptionId`, a `Callback` URL, and a `Filter` set to `oranapi.AlarmSubscriptionFilterNEW`.
2. The current time is saved before sending alerts.
3. One hundred alerts are concurrently sent to the `Alertmanager API client` using the `concurrenltySendAlerts` helper function.
4. The test waits for all notifications corresponding to the sent alerts to be received by the subscriber within a 20-minute timeout.
5. It is verified that the `sentAlerts` map is empty, indicating all alerts were received.
6. The test subscription is deleted.
