package tsparams

import "github.com/golang/glog"

const (
	// LabelSuite is the label for all tests in the PTP suite.
	LabelSuite = "ptp"

	// CloudEventsNamespace is the namespace used for the cloud events consumer and associated resources.
	CloudEventsNamespace = "cloud-events"
)

// LogLevel is the glog level used for all helpers in the PTP suite. It is set so that eco-goinfra is 100, cnf/ran is
// 90, and the suite itself is 80.
const LogLevel glog.Level = 80
