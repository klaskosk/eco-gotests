# Test Case Summary for 83556

Test case 83556 is located in tests/cnf/ran/oran/tests/oran-alarms.go and is named "acknowledges an alarm and listen for notification".

## Goal

The goal of this test case is to verify that an alarm can be acknowledged via the O2IMS API and that a corresponding notification is received by a subscriber.

## Test Setup

Prior to the test case, an `O2IMS API client` and an `Alertmanager API client` are created. The `spoke 1 cluster ID` is also retrieved. These are done in the `BeforeEach` block.

It does not require a git config set up such that X.

## Test Steps

1. A test alarm with `SeverityMajor` is created and sent to the `Alertmanager API client`.
2. The test waits for the created alarm to exist in the `O2IMS API`.
3. A test subscription is created with a `ConsumerSubscriptionId` and a `Callback` URL.
4. The current time is saved before acknowledging the alarm.
5. The alarm is acknowledged by patching its `AlarmAcknowledged` status to `true` via the `O2IMS API`.
6. The test waits for a notification to be received by the subscriber, matching the `tracker` of the sent alarm and verifying the `NotificationEventType` is `ACKNOWLEDGE`.
7. The test subscription is deleted.
