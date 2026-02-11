package stability

import (
	"fmt"
	"strings"

	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/daemonlogs"
)

// DefaultOffsetThresholdAbsoluteNanoseconds is the default absolute offset threshold used by stability analysis. It is
// intentionally positive, since it is meant for comparison with absolute offsets.
const DefaultOffsetThresholdAbsoluteNanoseconds int64 = 100

// OffsetStatistics captures descriptive statistics about absolute offsets over a set of samples. All offsets are in
// nanoseconds.
type OffsetStatistics struct {
	// MaxAbs is the maximum absolute offset in nanoseconds.
	MaxAbs int64
	// MinAbs is the minimum absolute offset in nanoseconds.
	MinAbs int64
	// AvgAbs is the average absolute offset in nanoseconds.
	AvgAbs float64
	// SampleCount is the number of samples used to compute the statistics.
	SampleCount int
}

// StateTransition describes a ptp4l state change between adjacent parsed entries.
type StateTransition struct {
	// From is the servo state of the previous entry. It is the letter s followed by a number. For example, "s1"
	// means the servo is in state 1.
	From string
	// To is the servo state of the current entry. It is the letter s followed by a number. For example, "s1" means
	// the servo is in state 1.
	To string
	// Raw is the full raw log line that was matched.
	Raw string
}

// AnalysisResult is the pass/fail decision output of stability analysis.
type AnalysisResult struct {
	// Passed is true if the analysis passed, false otherwise. It counts as passed if there are no failure details.
	Passed bool
	// Details is a list of failure detail messages. It is empty if the analysis passed.
	Details []string

	// PTP4LStats is the descriptive statistics for the ptp4l process.
	PTP4LStats OffsetStatistics
	// PHC2SYSStats is the descriptive statistics for the phc2sys process.
	PHC2SYSStats OffsetStatistics

	// PTP4LStartCount is the number of times the ptp4l process was started.
	PTP4LStartCount uint

	// FaultyLines is a list of lines containing the word "FAULTY".
	FaultyLines []string
	// TimeoutLines is a list of lines containing the word "timeout".
	TimeoutLines []string
	// StateTransitions is a list of servo state transitions between adjacent entries.
	StateTransitions []StateTransition

	// PTP4LThresholdViolations is a list of ptp4l log entries whose absolute offset in nanoseconds exceeds the
	// threshold.
	PTP4LThresholdViolations []daemonlogs.LogEntry
	// PHC2SYSThresholdViolations is a list of phc2sys log entries whose absolute offset in nanoseconds exceeds the
	// threshold.
	PHC2SYSThresholdViolations []daemonlogs.LogEntry

	// ParseWarnings is a list of warnings that occurred during parsing. These are warnings where the log lines
	// could not be parsed into log entries.
	ParseWarnings []string
}

// Analyze evaluates parsed daemon logs against stability policy.
func Analyze(parsed daemonlogs.ParsedLogs, thresholdAbsoluteNanoseconds int64) AnalysisResult {
	if thresholdAbsoluteNanoseconds <= 0 {
		thresholdAbsoluteNanoseconds = DefaultOffsetThresholdAbsoluteNanoseconds
	}

	result := AnalysisResult{
		PTP4LStats:      calculateStatisticsFromEntries(parsed.PTP4L.Entries),
		PHC2SYSStats:    calculateStatisticsFromEntries(parsed.PHC2SYS.Entries),
		PTP4LStartCount: parsed.PTP4LStartCount,
		ParseWarnings:   buildParseWarnings(parsed),
	}

	result.FaultyLines = collectLinesContaining(parsed.Lines, "faulty")
	result.TimeoutLines = collectLinesContaining(parsed.Lines, "timeout")
	result.PTP4LThresholdViolations = findThresholdViolations(parsed.PTP4L.Entries, thresholdAbsoluteNanoseconds)
	result.PHC2SYSThresholdViolations = findThresholdViolations(parsed.PHC2SYS.Entries, thresholdAbsoluteNanoseconds)
	result.StateTransitions = findStateTransitions(parsed.PTP4L.Entries)

	result.Details = buildFailureDetails(result)
	result.Passed = len(result.Details) == 0

	return result
}

// DiagnosticMessage builds a concise multi-line summary for assertions and report entries.
func (result AnalysisResult) DiagnosticMessage() string {
	var builder strings.Builder

	if len(result.Details) == 0 {
		builder.WriteString("No stability anomalies detected.")
	} else {
		builder.WriteString("Stability anomalies detected:")

		for _, detail := range result.Details {
			builder.WriteString("\n- ")
			builder.WriteString(detail)
		}
	}

	if len(result.ParseWarnings) > 0 {
		builder.WriteString("\nParse warnings:")

		for _, warning := range result.ParseWarnings {
			builder.WriteString("\n- ")
			builder.WriteString(warning)
		}
	}

	builder.WriteByte('\n')
	builder.WriteString(formatStatsLine("ptp4l", result.PTP4LStats))
	builder.WriteByte('\n')
	builder.WriteString(formatStatsLine("phc2sys", result.PHC2SYSStats))
	fmt.Fprintf(&builder, "\nptp4l_start_count=%d", result.PTP4LStartCount)

	return builder.String()
}

// buildParseWarnings returns human-readable warnings for any lines that were dropped during parsing.
func buildParseWarnings(parsed daemonlogs.ParsedLogs) []string {
	var warnings []string

	if parsed.PTP4L.DroppedLines > 0 {
		warnings = append(warnings, fmt.Sprintf(
			"ptp4l dropped %d/%d candidate delay lines during parsing",
			parsed.PTP4L.DroppedLines, parsed.PTP4L.CandidateLines))
	}

	if parsed.PHC2SYS.DroppedLines > 0 {
		warnings = append(warnings, fmt.Sprintf(
			"phc2sys dropped %d/%d candidate delay lines during parsing",
			parsed.PHC2SYS.DroppedLines, parsed.PHC2SYS.CandidateLines))
	}

	return warnings
}

// buildFailureDetails collects one detail string per stability check that failed. These are global checks based on the
// entire analysis results.
func buildFailureDetails(result AnalysisResult) []string {
	var details []string

	if result.PTP4LStats.SampleCount == 0 {
		details = append(details, "no ptp4l delay logs parsed")
	}

	if result.PHC2SYSStats.SampleCount == 0 {
		details = append(details, "no phc2sys delay logs parsed")
	}

	if len(result.FaultyLines) > 0 {
		details = append(details, fmt.Sprintf("found %d lines containing FAULTY", len(result.FaultyLines)))
	}

	if len(result.TimeoutLines) > 0 {
		details = append(details, fmt.Sprintf("found %d lines containing timeout", len(result.TimeoutLines)))
	}

	if len(result.PTP4LThresholdViolations) > 0 {
		details = append(details,
			fmt.Sprintf("found %d ptp4l s2 offset violations over threshold", len(result.PTP4LThresholdViolations)))
	}

	if len(result.PHC2SYSThresholdViolations) > 0 {
		details = append(details, fmt.Sprintf(
			"found %d phc2sys s2 offset violations over threshold", len(result.PHC2SYSThresholdViolations)))
	}

	if len(result.StateTransitions) > 0 {
		details = append(details, fmt.Sprintf("found %d ptp4l state transitions", len(result.StateTransitions)))
	}

	if result.PTP4LStartCount > 1 {
		details = append(details, fmt.Sprintf("found %d ptp4l restarts", result.PTP4LStartCount-1))
	}

	return details
}

// collectLinesContaining returns all lines that contain needle (case-insensitive).
func collectLinesContaining(lines []string, needle string) []string {
	var matches []string

	lowerNeedle := strings.ToLower(needle)

	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), lowerNeedle) {
			matches = append(matches, line)
		}
	}

	return matches
}

// findThresholdViolations returns s2 entries whose absolute offset in nanoseconds exceeds thresholdAbsoluteNanoseconds.
func findThresholdViolations(entries []daemonlogs.LogEntry, thresholdAbsoluteNanoseconds int64) []daemonlogs.LogEntry {
	var violations []daemonlogs.LogEntry

	for _, entry := range entries {
		if entry.State != "s2" {
			continue
		}

		if abs(entry.Offset) > thresholdAbsoluteNanoseconds {
			violations = append(violations, entry)
		}
	}

	return violations
}

// findStateTransitions returns transitions where the state changed between adjacent entries.
func findStateTransitions(entries []daemonlogs.LogEntry) []StateTransition {
	var transitions []StateTransition

	var previousState string

	for _, entry := range entries {
		if previousState != "" && previousState != entry.State {
			transitions = append(transitions, StateTransition{
				From: previousState,
				To:   entry.State,
				Raw:  entry.Raw,
			})
		}

		previousState = entry.State
	}

	return transitions
}

// calculateStatisticsFromEntries extracts offsets from log entries and computes aggregate statistics.
func calculateStatisticsFromEntries(entries []daemonlogs.LogEntry) OffsetStatistics {
	values := make([]int64, 0, len(entries))
	for _, entry := range entries {
		values = append(values, entry.Offset)
	}

	return calculateStatistics(values)
}

// calculateStatistics computes max, min, and average absolute offset over the given values.
func calculateStatistics(values []int64) OffsetStatistics {
	stats := OffsetStatistics{
		SampleCount: len(values),
	}

	if len(values) == 0 {
		return stats
	}

	firstAbs := abs(values[0])
	minAbs := firstAbs
	maxAbs := firstAbs
	totalAbs := float64(firstAbs)

	for _, value := range values[1:] {
		absoluteValue := abs(value)
		totalAbs += float64(absoluteValue)

		if absoluteValue < minAbs {
			minAbs = absoluteValue
		}

		if absoluteValue > maxAbs {
			maxAbs = absoluteValue
		}
	}

	stats.AvgAbs = totalAbs / float64(len(values))
	stats.MaxAbs = maxAbs
	stats.MinAbs = minAbs

	return stats
}

// formatStatsLine renders an OffsetStatistics value as a single key=value diagnostic line.
func formatStatsLine(process string, stats OffsetStatistics) string {
	return fmt.Sprintf("%s_offsets_max_abs=%d min_abs=%d avg_abs=%.3f samples=%d",
		process, stats.MaxAbs, stats.MinAbs, stats.AvgAbs, stats.SampleCount)
}

// abs returns the absolute value of an int64. It ignores the possibility of overflow since it is not applicable to the
// use case of PTP offsets, which should never be the minimum int64 value.
func abs(value int64) int64 {
	if value < 0 {
		return -value
	}

	return value
}
