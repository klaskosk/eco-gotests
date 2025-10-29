# Test Case Summary for 82919

Test case 82919 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verifies connectivity between pods from statefuleset scheduled on the same node post hard reboot".

## Goal

The goal of this test case is to verify network connectivity between pods belonging to a statefulset, scheduled on the same node, after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that a statefulset with pods deployed on the same node is present and has recovered its state.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.EnsurePodConnectivityOnSameNodeAfterNodePowerOff`) to validate the connectivity between the statefulset pods on the same node.
2. This typically involves performing network connectivity tests (e.g., ping, curl) between the pods within the statefulset.
3. The intent is to confirm that the statefulset pods can re-establish communication and maintain their network identity on the same node after a disruptive cluster event.
