# Test Case Summary for Successful Alarm Retrieval

Test case "Successful alarm retrieval from the API" is located in tests/system-tests/o-cloud/tests/alarms.go and is named "It verifies successful alarm retrieval from the API".

## Goal

The goal of this test case is to verify the successful retrieval of an alarm from the API.

## Test Setup

Prior to the test case, the following changes are needed:

- Deploy the subscriber for alarm notifications.
- Ensure BMHs (OCloudConfig.BmhSpoke1, OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.

It does not require a git config set up.

## Test Steps

1. Verify that both BMHs (OCloudConfig.BmhSpoke1 and OCloudConfig.BmhSpoke2) are available in OCloudConfig.InventoryPoolNamespace.
2. Provision a SNO cluster using `VerifyProvisionSnoCluster` with `OCloudConfig.TemplateName`, `OCloudConfig.TemplateVersionAISuccess`, `OCloudConfig.NodeClusterName1`, `OCloudConfig.OCloudSiteID`, `ocloudparams.PolicyTemplateParameters`, and `ocloudparams.ClusterInstanceParameters1`.
3. Verify that the OCloud Custom Resources (CRs) for the provisioned cluster exist using `VerifyOcloudCRsExist`.
4. Verify that the cluster instance creation is completed using `VerifyClusterInstanceCompleted`.
5. Verify that all policies in the cluster's namespace are compliant using `VerifyAllPoliciesInNamespaceAreCompliant`.
6. Verify that the provisioning request is fulfilled using `VerifyProvisioningRequestIsFulfilled`.
7. Create an O2IMS client using `createO2IMSClient` and an alarm subscription using `createAlarmSubscription`.
8. Modify the PTP operator deployment resources to trigger an alarm:
    a. Create a SNO API client for `OCloudConfig.ClusterName1` using `CreateSnoAPIClient`.
    b. Get the current time from the hub using `getHubCurrentTime`.
    c. Verify all pods are running in the PTP namespace using `VerifyAllPodsRunningInNamespace`.
    d. Modify PTP operator resources to trigger an alarm using `modifyPTPOperatorResources(snoAPIClient, true)`.
    e. Verify policies are not compliant using `VerifyPoliciesAreNotCompliant`.
    f. Wait for `AlarmWaitTime`.
9. Filter alarms by alertname since the `startTime` using `getACMPolicyViolationAlarmsSinceStartTime` and verify a minimum of 1 alarm is found using `verifyMinimumAlarmCount`.
10. For each filtered alarm, verify its successful retrieval using `alarmsClient.GetAlarm(alarm.AlarmEventRecordId)`.
11. Modify the PTP operator deployment resources to stop triggering the alarm using `modifyPTPOperatorResources(snoAPIClient, false)`.
12. Verify that all policies in the cluster's namespace are compliant using `VerifyAllPoliciesInNamespaceAreCompliant`.
13. Clean up the alarm subscription using `cleanupAlarmSubscription`.
14. Store the `provisioningRequest` and `clusterInstance` in `sharedProvisioningRequest` and `sharedClusterInstance` respectively.
