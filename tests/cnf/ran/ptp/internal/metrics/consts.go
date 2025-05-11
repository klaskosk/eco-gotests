package metrics

// InterfaceClockRealtime is the name of the interface representing the realtime clock. It is a string since there is no
// enum for interfaces.
const InterfaceClockRealtime = "CLOCK_REALTIME"

// InterfaceMaster is the name of the interface representing the master clock. It is a string since there is no enum for
// interfaces.
const InterfaceMaster = "master"

// PtpMetricKey is an enum representing all possible keys for labels on PTP metrics.
type PtpMetricKey string

//nolint:revive // The key names are self explanatory and do not need individual comments.
const (
	KeyProcess   PtpMetricKey = "process"
	KeyInterface PtpMetricKey = "iface"
	KeyNode      PtpMetricKey = "node"
	KeyConfig    PtpMetricKey = "config"
	KeyProfile   PtpMetricKey = "profile"
	KeyThreshold PtpMetricKey = "threshold"
	KeyFrom      PtpMetricKey = "from"
)

// PtpMetric is an enum representing all PTP metrics supported as typed queries.
type PtpMetric string

//nolint:revive // The metric names are self explanatory and do not need individual comments.
const (
	MetricClockState      PtpMetric = "openshift_ptp_clock_state"
	MetricProcessStatus   PtpMetric = "openshift_ptp_process_status"
	MetricThreshold       PtpMetric = "openshift_ptp_threshold"
	MetricNMEAStatus      PtpMetric = "openshift_ptp_nmea_status"
	MetricHAProfileStatus PtpMetric = "openshift_ptp_ha_profile_status"
	MetricPPSStatus       PtpMetric = "openshift_ptp_pps_status"
	MetricClockClass      PtpMetric = "openshift_ptp_clock_class"
)

// PtpClockState is an enum representing all possible states of the PTP clock.
type PtpClockState int

const (
	// ClockStateFreerun is the state of the PTP clock when it is not synchronized to a time transmitter.
	ClockStateFreerun PtpClockState = iota
	// ClockStateLocked is the state of the PTP clock when it is synchronized to a time transmitter.
	ClockStateLocked
	// ClockStateHoldover is the state of the PTP clock when it is in holdover mode, meaning it is using its
	// internal clock to maintain time when it is not receiving a signal from a time transmitter.
	ClockStateHoldover
)

// PtpProcessStatus is an enum representing all possible states of the PTP process.
type PtpProcessStatus int

//nolint:revive // The process status names are self explanatory and do not need individual comments.
const (
	ProcessStatusDown PtpProcessStatus = iota
	ProcessStatusUp
)

// PtpThresholdType is an enum representing all possible types of PTP thresholds. It corresponds to the keys of
// ptpv1.PtpClockThresholds.
type PtpThresholdType string

//nolint:revive // The threshold type names are self explanatory and do not need individual comments.
const (
	ThresholdHoldoverTimeout PtpThresholdType = "HoldOverTimeout"
	ThresholdMaxOffset       PtpThresholdType = "MaxOffsetThreshold"
	ThresholdMinOffset       PtpThresholdType = "MinOffsetThreshold"
)

// PtpNMEAStatus is an enum representing all possible states of the PTP NMEA status. It is similar to PtpProcessStatus
// but typed specifically for the NMEA status metric.
type PtpNMEAStatus int

//nolint:revive // The NMEA status names are self explanatory and do not need individual comments.
const (
	NMEAStatusUnavailable PtpNMEAStatus = iota
	NMEAStatusAvailable
)

// PtpHAProfileStatus is an enum representing all possible states of the PTP HA profile status. It is similar to
// PtpProcessStatus but typed specifically for the HA profile status metric.
type PtpHAProfileStatus int

//nolint:revive // The HA profile status names are self explanatory and do not need individual comments.
const (
	HAProfileStatusInactive PtpHAProfileStatus = iota
	HAProfileStatusActive
)

// PtpPPSStatus is an enum representing all possible states of the PTP PPS status. It is similar to PtpProcessStatus but
// typed specifically for the PPS status metric.
type PtpPPSStatus int

//nolint:revive // The PPS status names are self explanatory and do not need individual comments.
const (
	PPSStatusUnavailable PtpPPSStatus = iota
	PPSStatusAvailable
)

// PtpProcess is an enum representing all possible values for PTP processes. This is used as the type for the from label
// and the process label.
type PtpProcess string

//nolint:revive // The process names are self explanatory and do not need individual comments.
const (
	ProcessPTP4L   PtpProcess = "ptp4l"
	ProcessPHC2SYS PtpProcess = "phc2sys"
	ProcessTS2PHC  PtpProcess = "ts2phc"
	ProcessGPSD    PtpProcess = "gpsd"
	ProcessGPSPIPE PtpProcess = "gpspipe"
	ProcessDPLL    PtpProcess = "dpll"
	ProcessGNSS    PtpProcess = "gnss"
	ProcessGM      PtpProcess = "GM"
)
