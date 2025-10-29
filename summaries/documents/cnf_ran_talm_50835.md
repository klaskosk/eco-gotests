# Test Case Summary for 50835

Test case 50835 is located in tests/cnf/ran/talm/tests/talm-backup.go and is named "It should have a failed cgu for single spoke".

## Goal

The goal of this test case is to verify that a ClusterGroupUpgrade (CGU) fails with an "UnrecoverableError" status when a single spoke cluster has insufficient disk space for a backup.

## Test Setup

Prior to the test case, the following changes are needed:

- The TALM version on the hub must be between 4.11 and 4.15 (exclusive of 4.16).
- Both hub and spoke1 API clients must be present.
- The filesystem on spoke1 is prepared to simulate low disk space using `mount.PrepareEnvWithSmallMountPoint`.

It does not require a git config set up such that X.

## Test Steps

1. A CGU is created with `cgu.NewCguBuilder`, targeting `RANConfig.Spoke1Name` and `tsparams.PolicyName`.
2. The `Backup` field in the CGU definition is set to `true`.
3. The CGU is set up using `helper.SetupCguWithNamespace`.
4. The test waits for the CGU backup status for `RANConfig.Spoke1Name` to become "UnrecoverableError" using `assertBackupStatus`.
