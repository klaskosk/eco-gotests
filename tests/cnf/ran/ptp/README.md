# PTP Test Cases

## `ptp_recovery.go`

* [x] **59862**: "should recover the phc2sys process after killing it"
* [x] **57197**: "should create a new ptp4l process after killing a ptp4l process that is not related to the phc2sy process"
* [ ] **49736**: "should restart both ptp4l processes with one related to phc2sys after killing them"
* [ ] **49737**: "should recover the ptp4l process after the killing a ptp4l process that is related to phc2sys process"
* [ ] **59863**: "should recover the ts2phc process after the killing a ts2phc process"
* [ ] **59864**: "should recover the ptp4l process after the killing a ptp4l process that is related to ts2phc process"
* [ ] **64777**: "should recover gpsd process after killing it on node"
* [x] **64775**: "should recover to stable state after delete PTP daemon pod"
* [x] **54245**: "validates HTTP PTP events via consumer"
  * SKIPPED: the test for 49741 already covers this
* [x] **59996**: "validates the system is fully functional after removing consumer"
  * SKIPPED: not a realistic customer use case
* [x] **82218**: "validates the consumer events after ptpoperatorconfig api version is modified"
* [x] **59858**: "should return to same stable status after ptp node soft reboot"
* [x] **59995**: "validates PTP consumer events after ptp node reboot"
* [ ] **70111**: "should make nmea lost after GPS cold reboot"
* [ ] **78463**: "verifies t-gm transition from holdover to locked due to gnss recovery"
* [ ] **78464**: "verifies t-gm transition from holdover to freerun due to timeout"
* [ ] **78465**: "verifies t-gm transition from holdover to freerun due to offset"
* [ ] **81205**: "checks FREERUN status are generated for dpll process for RX interface and GM process for TX interface"

## `ptp_interfaces.go`

* [x] **49742**: "should generate events when slave interface goes down and up"
* [ ] **49734**: "should have no effect when Boundary Clock master interface goes down and up"
* [ ] **73093**: "should change high availability active profile when other nic interface is down"
* [ ] **73094**: "should move to FREERUN state when active and inactive interfaces are down"
* [ ] **73095**: "should change high availability active profile when active profile is deleted"
* [ ] **80963**: "verifies 2-port oc ha failover when active port goes down"
* [ ] **80964**: "verifies 2-port oc ha holdover & freerun when both ports go down"
* [ ] **82012**: "verifies 2-port oc ha passive interface recovery"

## `ptp_events_and_metrics.go`

* [x] **82480**: "should have [LOCKED] clock state in PTP metrics"
* [x] **66848**: "should have the 'phc2sys' and 'ptp4l' processes in 'UP' state in PTP metrics"
* [x] **49741**: "should change the slave clock state to free run after modify the offset threshold"
* [x] **82302**: "should have the 'phc2sys' and 'ptp4l' processes 'UP' after ptp config change"

## `z_ptp_leap_file.go`

* [ ] **75325**: "should add leap event announcement in leap configmap when removing the last announcement"

---

# Recent Commits Affecting the Test Directory

## 1. **Add HOLDOVER state move to LOCKED check without getting FREERUN event**

**Commit:** `321474ee`
**Author:** Hen Shay Hassid
**Date:** Mon Jun 23 17:00:27 2025 +0300

**Files Modified:**
* `test/ran/ptp/ptp_suite_test.go` (12 lines)
* `test/ran/ptp/ranptphelper/ranptphelper.go` (65 lines)
* `test/ran/ptp/tests/ptp_events_and_metrics.go` (5 lines)
* `test/ran/ptp/tests/ptp_recovery.go` (72 lines)

**Key Changes:**
* **Enhanced PTP Configuration Management:** Updated `UpdatePtpConfigSpecs` function to accept PTP daemon pods and compare configurations before applying changes, avoiding unnecessary updates
* **Added HOLDOVER State Handling:** Introduced version-specific logic for PTP 4.20+ to expect HOLDOVER events instead of FREERUN events when ptp4l processes are killed
* **New Helper Functions:** Added `WaitForPtpConfigToBeLoaded()` and `ComparePtpSpecs()` for better PTP configuration validation
* **Improved Recovery Logic:** Added `increaseHoldoverTimeout()` function to set HoldOverTimeout to 180 seconds for testing
* **Event Validation:** Added validation to ensure no FREERUN events are received when expected HOLDOVER events occur

## 3. **changed the accepted version from 4.19 to 4.18 in OC 2 port test cases to accommodate for backporting**

**Commit:** `755979d1`
**Author:** Daniel Popsuevich
**Date:** Mon Jul 7 17:29:31 2025 +0300

**Files Modified:**
* `test/ran/ptp/tests/ptp_interfaces.go` (1 line)

**Key Changes:**
* **Version Compatibility:** Changed minimum PTP version requirement from 4.19 to 4.18 for OC 2 port interface tests to support backporting scenarios

## 4. **OC 2 port recovery and validation**

**Commit:** `cbdd3f69`
**Author:** Daniel Popsuevich
**Date:** Tue Jun 10 17:23:07 2025 +0000

**Files Modified:**
* `test/ran/ptp/tests/ptp_interfaces.go` (29 lines)

**Key Changes:**
* **Simplified Test Structure:** Removed redundant preflight and post-test validation logic
* **New Helper Function:** Added `restoreValidateOc2Port()` function to restore interfaces and validate OC 2 port states
* **Enhanced Recovery Testing:** Added automatic restoration and validation after interface down/up scenarios
* **Improved Test Flow:** Streamlined test execution by removing duplicate validation steps
