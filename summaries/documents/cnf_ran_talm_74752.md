# Test Case Summary for 74752

Test case 74752 is located in tests/cnf/ran/talm/tests/talm-backup.go and is named "It should not affect backup on second spoke in same batch".

## Goal

The goal of this test case is to verify that when two spoke clusters are in the same batch for a CGU with backup enabled, and one spoke fails due to insufficient disk space, the backup on the second spoke is not affected and succeeds.

## Test Setup

Prior to the test case, the following changes are needed:

- The TALM version on the hub must be between 4.11 and 4.15 (exclusive of 4.16).
- The hub, spoke1, and spoke2 API clients must be present.
- The filesystem on spoke1 is prepared to simulate low disk space using `mount.PrepareEnvWithSmallMountPoint`.

It does not require a git config set up such that X.

## Test Steps

1. A CGU is created with `cgu.NewCguBuilder`, with a `maxConcurrency` of 2, targeting `RANConfig.Spoke1Name` and `RANConfig.Spoke2Name` and `tsparams.PolicyName`. This ensures both spokes are in the same batch.
2. The `Backup` field in the CGU definition is set to `true`.
3. The CGU is set up using `helper.SetupCguWithNamespace`.
4. The test waits for the CGU backup status for `RANConfig.Spoke1Name` to become "UnrecoverableError" using `assertBackupStatus`.
5. The test then waits for the CGU backup status for `RANConfig.Spoke2Name` to become "Succeeded" using `assertBackupStatus`, verifying that the failure of spoke1 did not affect spoke2.
