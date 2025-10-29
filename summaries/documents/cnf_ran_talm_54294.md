# Test Case Summary for 54294

Test case 54294 is located in tests/cnf/ran/talm/tests/talm-backup.go and is named "It verifies backup begins and succeeds after CGU is enabled".

## Goal

The goal of this test case is to verify that when a ClusterGroupUpgrade (CGU) is initially disabled with backup enabled, the backup process does not start until the CGU is enabled, after which it successfully completes.

## Test Setup

Prior to the test case, the following changes are needed:

- The TALM version on the hub must be at least 4.12.
- Both hub and spoke1 API clients must be present.

It does not require a git config set up such that X.

## Test Steps

1. A CGU is created with `cgu.NewCguBuilder`, targeting `RANConfig.Spoke1Name` and `tsparams.PolicyName`.
2. The `Backup` field in the CGU definition is set to `true`, and the `Enable` field is set to `false` using `ptr.To(false)`. The `Timeout` for remediation strategy is set to 30.
3. The CGU is set up using `helper.SetupCguWithNamespace`.
4. The test verifies that backup *does not* begin when the CGU is disabled by calling `cguBuilder.WaitUntilBackupStarts` and expecting an error.
5. The CGU is then enabled by setting `cguBuilder.Definition.Spec.Enable = ptr.To(true)` and updating the CGU.
6. The test then waits for backup to begin using `cguBuilder.WaitUntilBackupStarts`.
7. Finally, the test waits for the CGU to indicate backup succeeded for the spoke using `assertBackupStatus`.
