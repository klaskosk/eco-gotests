# Test Case Summary for 47948

Test case 47948 is located in tests/cnf/ran/talm/tests/talm-blockingcr.go and is named "verifies CGU succeeded with blocking CR".

## Goal

The goal of this test case is to verify that a ClusterGroupUpgrade (CGU) eventually succeeds when it is configured with a blocking CR, and that the blocked CGU only proceeds after the blocking CGU has completed successfully.

## Test Setup

Prior to the test case, the following changes are needed:

- The TALM operator version on the hub cluster must be 4.12 or higher.

After each test, the following cleanup actions are performed:

- Cleanup test resources associated with blocking CRs A and B on the hub cluster.
- Delete temporary namespaces (`tsparams.TemporaryNamespace + blockingA` and `tsparams.TemporaryNamespace + blockingB`) on spoke 1.

It does not require a git config set up such that X.

## Test Steps

1.  Two ClusterGroupUpgrades (CGU) are defined: `cguA` with a 10-minute timeout and `cguB` with a 15-minute timeout. `cguB` is configured to be blocked by `cguA` by setting its `Spec.BlockingCRs` to reference `cguA`.
2.  Both `cguA` and `cguB` are set up using `helper.SetupCguWithNamespace`. This involves creating a temporary namespace, creating a policy for that namespace, creating policy components (PolicySet, PlacementRule, PlacementBinding), waiting for these components to exist, and then creating the CGU itself on the hub cluster.
3.  After a system stabilization period, both `cguA` and `cguB` are enabled.
4.  The test verifies that `cguB` is blocked by `cguA` by waiting for the `Progressing` condition of `cguB` to be `False` with a specific message indicating it's blocked by `cguA`.
5.  The test then waits for `cguA` to reach a successful finish condition.
6.  Finally, the test waits for `cguB` to also reach a successful finish condition, demonstrating that it proceeded only after `cguA` completed.
