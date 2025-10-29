# Test Case Summary for Successful Alarms Cleanup

Test case "Successful alarms cleanup from the database after the retention period" is located in tests/system-tests/o-cloud/tests/alarms.go and is named "It verifies successful alarms cleanup from the database after the retention period".

## Goal

The goal of this test case is to verify that alarms are successfully cleaned up from the database after the configured retention period.

## Test Setup

Prior to the test case, the following changes are needed:

- Deploy the subscriber for alarm notifications.

It does not require a git config set up.

## Test Steps

1. Patch the alarm service configuration:
    a. Create an O2IMS client using `createO2IMSClient`.
    b. Create an alarm subscription using `createAlarmSubscription`.
    c. Create a SNO API client for `OCloudConfig.ClusterName1` using `CreateSnoAPIClient`.
    d. Get the current time from the hub using `getHubCurrentTime`.
    e. Patch the alarm service configuration to set the `RetentionPeriod` to 1 day using `alarmsClient.PatchAlarmServiceConfiguration` with `DefaultRetentionPeriod`.
    f. Verify that the patch was successful and the `RetentionPeriod` is set to 1.
2. Loop `ExpectedAlarmCount` times to trigger and clear alarms:
    a. Verify all pods are running in the PTP namespace using `VerifyAllPodsRunningInNamespace`.
    b. Modify PTP operator deployment resources to trigger an alarm using `modifyPTPOperatorResources(snoAPIClient, true)`.
    c. Verify policies are not compliant using `VerifyPoliciesAreNotCompliant`.
    d. Wait for `AlarmWaitTime`.
    e. Modify PTP operator resources to stop triggering the alarm using `modifyPTPOperatorResources(snoAPIClient, false)`.
    f. Verify policies are compliant using `VerifyAllPoliciesInNamespaceAreCompliant`.
    g. Wait for `AlarmWaitTime`.
3. Filter alarms and verify cleanup:
    a. Get the initial system time (`time.Now()`).
    b. Calculate the `retentionPeriod` based on `RetentionPeriodHours`.
    c. Filter alarms by alertname since `alarmsStartTime` using `getACMPolicyViolationAlarmsSinceStartTime`.
    d. Verify that at least `ExpectedAlarmCount` alarms exist using `verifyMinimumAlarmCount`.
    e. Wait until the `retentionPeriod` has passed to implicitly verify that alarms are cleaned up.
