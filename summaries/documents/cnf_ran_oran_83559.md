# Test Case Summary for 83559

Test case 83559 is located in tests/cnf/ran/oran/tests/oran-alarms.go and is named "updates alarm service configuration".

## Goal

The goal of this test case is to verify that the alarm service configuration, specifically the retention period, can be updated via the O2IMS API.

## Test Setup

Prior to the test case, an `O2IMS API client` and an `Alertmanager API client` are created. The `spoke 1 cluster ID` is also retrieved. These are done in the `BeforeEach` block.

It does not require a git config set up such that X.

## Test Steps

1. The current alarm service configuration is retrieved from the `O2IMS API`.
2. It is verified that the `originalRetentionPeriod` is at least 1 day.
3. The retention period is incremented by 1 and the configuration is patched via the `O2IMS API`.
4. It is verified that the `retention period` was indeed incremented by 1.
5. The retention period is decremented back to its original value and the configuration is updated via the `O2IMS API` using a PUT request.
6. It is verified that the `retention period` matches the original value.
