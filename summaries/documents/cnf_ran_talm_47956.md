# Test Case Summary for 47956

Test case 47956 is located in tests/cnf/ran/talm/tests/talm-blockingcr.go and is named "verifies CGU is blocked until blocking CR created and succeeded".

## Goal

The goal of this test case is to verify that a ClusterGroupUpgrade (CGU) remains blocked when its blocking CR is missing, and then proceeds to succeed once the blocking CR is created and successfully completes.

## Test Setup

Prior to the test case, the following changes are needed:

- The TALM operator version on the hub cluster must be 4.12 or higher.

After each test, the following cleanup actions are performed:

- Cleanup test resources associated with blocking CRs A and B on the hub cluster.
- Delete temporary namespaces (`tsparams.TemporaryNamespace + blockingA` and `tsparams.TemporaryNamespace + blockingB`) on spoke 1.

It does not require a git config set up such that X.

## Test Steps

1.  Two ClusterGroupUpgrades (CGU) are defined: `cguA` with a 10-minute timeout and `cguB` with a 15-minute timeout. `cguB` is configured to be blocked by `cguA` by setting its `Spec.BlockingCRs` to reference `cguA`.
2.  `cguB` is set up using `helper.SetupCguWithNamespace`. This involves creating a temporary namespace, creating a policy for that namespace, creating policy components (PolicySet, PlacementRule, PlacementBinding), waiting for these components to exist, and then creating the CGU itself on the hub cluster.
3.  After a system stabilization period, `cguB` is enabled.
4.  The test verifies that `cguB` is blocked because `cguA` is missing, by waiting for the `Progressing` condition of `cguB` to be `False` with a specific message indicating it's blocked by a missing `cguA`.
5.  `cguA` is then set up using `helper.SetupCguWithNamespace` and subsequently enabled.
6.  The test waits for `cguA` to reach a successful finish condition.
7.  Finally, the test waits for `cguB` to also reach a successful finish condition, demonstrating that it proceeded only after `cguA` was created and completed successfully.
