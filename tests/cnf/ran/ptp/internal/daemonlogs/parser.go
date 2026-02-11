package daemonlogs

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/processes"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
	"k8s.io/klog/v2"
)

// LogEntry is a parsed delay log line from a synchronization daemon (ptp4l or phc2sys). Each process contains a clock
// offset and a servo state, both of which are extracted from the log line.
type LogEntry struct {
	// Raw is the full raw log line that was matched.
	Raw string
	// Offset is the offset of the log line in nanoseconds.
	Offset int64
	// State is the servo state of the log line. It is the letter s followed by a number. For example, "s1" means
	// the servo is in state 1.
	State string
}

// ParseResult contains parsed log entries and the number of dropped candidate lines.
type ParseResult struct {
	// Entries is the parsed log entries from log lines that match the pattern.
	Entries []LogEntry
	// DroppedLines is the number of candidate log lines that were dropped because the integers could not be parsed.
	DroppedLines int
	// CandidateLines is the number of log lines that matched the log pattern.
	CandidateLines int
}

// ParsedLogs contains parsed daemon streams and full raw lines.
type ParsedLogs struct {
	// Lines is the full raw log lines that were collected.
	Lines []string
	// PTP4L is the parsed log entries from log lines that matched the ptp4l pattern.
	PTP4L ParseResult
	// PHC2SYS is the parsed log entries from log lines that matched the phc2sys pattern.
	PHC2SYS ParseResult
	// PTP4LStartCount is the number of times the ptp4l process was started.
	PTP4LStartCount uint
}

var (
	// ptp4lPattern is a regular expression that matches the ptp4l log lines. For example:
	//  ptp4l[401304.873]: [ptp4l.1.config:6] master offset         -3 s2 freq  -94379 path delay       161
	ptp4lPattern = regexp.MustCompile(`^ptp4l\[.*?\boffset\s+(?P<offset>-?\d+)\s+(?P<state>s\d+).*delay`)
	// phc2sysPattern is a regular expression that matches the phc2sys log lines. For example:
	//  phc2sys[401304.879]: [ptp4l.1.config:6] CLOCK_REALTIME phc offset        -5 s2 freq  -19334 delay    470
	phc2sysPattern = regexp.MustCompile(`^phc2sys\[.*?\boffset\s+(?P<offset>-?\d+)\s+(?P<state>s\d+).*delay`)
)

// ParseLogs parses raw daemon log lines into per-process streams.
func ParseLogs(lines []string) ParsedLogs {
	return ParsedLogs{
		Lines:           slices.Clone(lines),
		PTP4L:           parseDaemonLogs(lines, ptp4lPattern, processes.Ptp4l),
		PHC2SYS:         parseDaemonLogs(lines, phc2sysPattern, processes.Phc2sys),
		PTP4LStartCount: countPTP4LStarts(lines),
	}
}

// parseDaemonLogs extracts delay log entries matching pattern from lines. The pattern must contain named capture groups
// "offset" and "state", corresponding to the clock offset and servo state, respectively. The processName is used only
// for diagnostic logging of dropped lines.
func parseDaemonLogs(lines []string, pattern *regexp.Regexp, processName processes.PtpProcess) ParseResult {
	offsetIndex := pattern.SubexpIndex("offset")
	stateIndex := pattern.SubexpIndex("state")

	result := ParseResult{}

	for _, line := range lines {
		match := pattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		result.CandidateLines++

		offset, err := strconv.ParseInt(match[offsetIndex], 10, 64)
		if err != nil {
			klog.V(tsparams.LogLevel).Infof("%s: dropping line with unparseable offset %q: %v", processName, line, err)

			result.DroppedLines++

			continue
		}

		result.Entries = append(result.Entries, LogEntry{
			Raw:    line,
			Offset: offset,
			State:  match[stateIndex],
		})
	}

	return result
}

// countPTP4LStarts counts how many times ptp4l was started by looking for "Starting ptp4l" in the log lines.
func countPTP4LStarts(lines []string) uint {
	var startCount uint

	for _, line := range lines {
		if strings.Contains(line, "Starting ptp4l") {
			startCount++
		}
	}

	return startCount
}
