# Test Case Summary for 72463

Test case 72463 is located in tests/cnf/ran/gitopsztp/tests/ztp-argocd-node-delete.go and is named "Delete and re-add a worker node from cluster - should delete a worker node from the cluster".

## Goal

The goal of this test case is to verify that a worker node can be successfully deleted from a cluster using Argo CD GitOps, and that the cluster remains healthy after the deletion.

## Test Setup

Prior to the test case, the original clusters app source is saved and then reset after the test. The test checks for ZTP version 4.14 or later, verifies that the cluster is SNO+1 (single control plane and single worker node), and ensures the 'worker' machine config pool is ready. It identifies the worker node name and its BareMetalHost (BMH) namespace.

It does require a git config set up such that the clusters app can be updated via GitOps.

## Test Steps

1. Update the Argo CD Git path to apply a `crAnnotation` for node deletion and wait for synchronization.
2. Wait for the `crAnnotation` (`tsparams.NodeDeletionCrAnnotation`) to be added to the BareMetalHost of the worker node.
3. Reset the clusters app Git path to the original and then update it again to apply the suppression for node deletion, waiting for synchronization.
4. Wait for the worker node's BareMetalHost to be deprovisioned.
5. Check that the cluster is healthy and stable.
