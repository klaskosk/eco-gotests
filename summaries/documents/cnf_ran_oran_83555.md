# Test Case Summary for 83555

Test case 83555 is located in tests/cnf/ran/oran/tests/oran-alarms.go and is named "filters alarms from the API".

## Goal

The goal of this test case is to verify that alarms can be filtered from the O2IMS API based on specific criteria.

## Test Setup

Prior to the test case, an `O2IMS API client` and an `Alertmanager API client` are created. The `spoke 1 cluster ID` is also retrieved. These are done in the `BeforeEach` block.

It does not require a git config set up such that X.

## Test Steps

1. Two test alarms are created and sent to the `Alertmanager API client`: one with `SeverityMajor` and another with `SeverityMinor`.
2. The test waits for both the major and minor alarms to exist in the `O2IMS API`.
3. Alarms are filtered from the `O2IMS API` using a combined filter for `major severity` and `resourceID` matching the `spoke1ClusterID`.
4. It is verified that the filtered results contain the major alarm but do not contain the minor alarm.
