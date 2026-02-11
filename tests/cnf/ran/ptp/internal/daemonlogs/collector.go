package daemonlogs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/clients"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/ptpdaemon"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
)

// CollectionResult is the output of a long-window daemon log collection.
type CollectionResult struct {
	// NodeName is the name of the node that the logs were collected from.
	NodeName string
	// StartedAt is the time the collection started.
	StartedAt time.Time
	// EndedAt is the time the collection ended.
	EndedAt time.Time
	// Lines is the full raw log lines that were collected.
	Lines []string
	// Errors is the errors that occurred while collecting the logs.
	Errors []error
}

// CollectDaemonLogs collects linuxptp daemon logs for a single node for the provided duration. It polls every 10
// seconds for the full duration. Some log lines may be duplicated, but they will never be skipped. The returned pointer
// is non-nil if and only if error is nil.
//
// In the returned CollectionResult, the StartedAt and EndedAt times are the time the collection process started and
// ended, not necessarily the time the first and last log lines were collected.
func CollectDaemonLogs(client *clients.Settings, nodeName string, duration time.Duration) (*CollectionResult, error) {
	if client == nil {
		return nil, fmt.Errorf("cannot collect daemon logs with nil client")
	}

	if nodeName == "" {
		return nil, fmt.Errorf("cannot collect daemon logs with empty node name")
	}

	if duration <= 0 {
		return nil, fmt.Errorf("cannot collect daemon logs with non-positive duration: %s", duration)
	}

	startTime := time.Now()
	lastFetchTime := startTime

	result := CollectionResult{
		NodeName:  nodeName,
		StartedAt: startTime,
	}

	var collectionErrors []error

	err := wait.PollUntilContextTimeout(
		context.TODO(), 10*time.Second, duration, true, func(ctx context.Context) (bool, error) {
			// We save the time of the next last fetch before getting the logs to avoid missing log entries,
			// although this introduces the potential for duplicates.
			localFetchTime := time.Now()

			lines, fetchErr := collectLinesSince(client, nodeName, lastFetchTime)
			if fetchErr != nil {
				klog.V(tsparams.LogLevel).Infof("Error collecting daemon logs from node %s: %v", nodeName, fetchErr)

				collectionErrors = append(collectionErrors, fetchErr)
			} else {
				lastFetchTime = localFetchTime

				result.Lines = append(result.Lines, lines...)
			}

			return false, nil
		})

	// We expect to see a context.DeadlineExceeded error since we always poll the full duration of the timeout. In
	// practice, this means this if branch will never be reached, since the closure always returns a nil error.
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return nil, fmt.Errorf("unexpected error collecting daemon logs on node %s: %w", nodeName, err)
	}

	result.EndedAt = time.Now()
	result.Errors = collectionErrors

	return &result, nil
}

// collectLinesSince fetches daemon log lines produced after lastFetchTime from the PTP daemon pod on the given node. It
// returns the lines and an error if the fetch failed.
func collectLinesSince(
	client *clients.Settings, nodeName string, lastFetchTime time.Time) ([]string, error) {
	daemonPod, err := ptpdaemon.GetPtpDaemonPodOnNode(client, nodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get PTP daemon pod on node %s: %w", nodeName, err)
	}

	logs, err := daemonPod.GetLogsWithOptions(&corev1.PodLogOptions{
		SinceTime: &metav1.Time{Time: lastFetchTime},
		Container: ranparam.PtpContainerName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get logs from node %s since %s: %w", nodeName, lastFetchTime, err)
	}

	logLines := splitAndTrimLogLines(string(logs))

	return logLines, nil
}

// splitAndTrimLogLines splits a raw log string on newlines and discards empty lines.
func splitAndTrimLogLines(logs string) []string {
	lines := strings.Split(logs, "\n")
	filteredLines := make([]string, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		filteredLines = append(filteredLines, line)
	}

	return filteredLines
}
