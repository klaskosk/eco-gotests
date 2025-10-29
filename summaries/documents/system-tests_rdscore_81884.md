# Test Case Summary for 81884

Test case 81884 is located in tests/system-tests/rdscore/tests/00_validate_top_level.go and is named "Verify cluster log forwarding to the Kafka broker post hard reboot".

## Goal

The goal of this test case is to verify that cluster logs are successfully forwarded to a Kafka broker after an ungraceful (hard) cluster reboot.

## Test Setup

Prior to the test case, this test assumes that an ungraceful cluster reboot has occurred, and that a logging solution is configured to forward logs to a Kafka broker, and that a Kafka broker is accessible within the cluster and has recovered its state.

It does not require a git config set up.

## Test Steps

1. The test uses a helper function (`rdscorecommon.VerifyLogForwardingToKafka`) to validate the log forwarding functionality.
2. This typically involves generating logs within the cluster (e.g., by creating and deleting pods, or triggering application errors).
3. The test then verifies that these generated logs appear in the configured Kafka topics, confirming that the log forwarding mechanism is operational and resilient post-reboot.
4. The intent is to ensure the continuous and reliable delivery of cluster logs to an external Kafka broker even after a disruptive cluster event.
