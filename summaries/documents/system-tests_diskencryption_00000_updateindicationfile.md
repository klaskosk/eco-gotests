# Test Case Summary for Verifies update indication with file (/etc/host-hw-Updating.flag)

Test case "Verifies update indication with file (/etc/host-hw-Updating.flag)" is located in tests/system-tests/diskencryption/tests/tpm2.go.

## Goal

The goal of this test case is to verify that if the file `/etc/host-hw-Updating.flag` is present, disk decryption will succeed even if Secure Boot is disabled.

## Test Setup

Prior to the test case, the following changes are needed:

- The cluster should be a Single Node OpenShift (SNO) cluster.
- Secure Boot must be enabled initially.
- The root disk must be encrypted with TPM2 using PCR 1 and 7.
- TPM max retries and lockout counter are configured in the `BeforeEach` block.
- The TTY console options must be configured on the kernel boot line (nomodeset console=tty0 console=ttyS0,115200n8).

It does not require a git config set up.

## Test Steps

1. Verify that the Root disk is encrypted with tpm2 with PCR 1 and 7 using `helper.GetClevisLuksListOutput()` and `helper.LuksListContainsPCR1And7()`.
2. Check if Secure Boot is enabled using `BMCClient.IsSecureBootEnabled()`.
3. Verify that the reserved slot is not present using `helper.GetClevisLuksListOutput()` and `helper.LuksListContainsReservedSlot()`.
4. Create the `/etc/host-hw-Updating.flag` file using `file.TouchFile()`.
5. Disable Secure Boot using `BMCClient.SecureBootDisable()`.
6. Restart the node gracefully using `cluster.SoftRebootSNO()`.
7. Wait for the "pcr-rebind-boot" log to appear, indicating disk decryption success, using `stdinmatcher.WaitForRegex()`.
8. Wait for the cluster to recover using `cluster.WaitForRecover()`.
9. Enable Secure Boot using `BMCClient.SecureBootEnable()`.
10. Restart the node gracefully using `cluster.SoftRebootSNO()`.
