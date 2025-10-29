# Test Case Summary for 83554

Test case 83554 is located in tests/cnf/ran/oran/tests/oran-alarms.go and is named "retrieves an alarm from the API".

## Goal

The goal of this test case is to verify that an alarm can be retrieved from the O2IMS API after it has been created.

## Test Setup

Prior to the test case, an `O2IMS API client` and an `Alertmanager API client` are created. The `spoke 1 cluster ID` is also retrieved. These are done in the `BeforeEach` block.

It does not require a git config set up such that X.

## Test Steps

1. A test alarm with `SeverityMajor` is created and sent to the `Alertmanager API client`.
2. The test waits for the created alarm to exist in the `O2IMS API`.
3. The alarm is retrieved from the `O2IMS API` using its `AlarmEventRecordId`.
4. It is verified that the retrieved alarm's `tracker` extension matches the `tracker` of the initially sent alert.
