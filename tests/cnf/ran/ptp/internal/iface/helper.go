package iface

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ptp/internal/ptpdaemon"
)

// GetNICDriver uses ethtool to retrieve the driver for a given network interface on a specified node.
func GetNICDriver(client *clients.Settings, nodeName string, ifName Name) (string, error) {
	command := fmt.Sprintf("ethtool -i %s | grep --color=no driver | awk '{print $2}'", ifName)
	output, err := ptpdaemon.ExecuteCommandInPtpDaemonPod(client, nodeName, command)

	if err != nil {
		return "", fmt.Errorf("failed to get NIC driver for interface %s on node %s: %w", ifName, nodeName, err)
	}

	return strings.TrimSpace(output), nil
}

// GetPTPHardwareClock uses ethtool to retrieve the PTP hardware clock for a given network interface on a specified
// node.
func GetPTPHardwareClock(client *clients.Settings, nodeName string, ifName Name) (int, error) {
	command := fmt.Sprintf("ethtool -T %s | grep 'PTP Hardware Clock' | cut -d' ' -f4", ifName)
	output, err := ptpdaemon.ExecuteCommandInPtpDaemonPod(client, nodeName, command)

	if err != nil {
		return -1, fmt.Errorf("failed to get PTP hardware clock for interface %s on node %s: %w", ifName, nodeName, err)
	}

	hardwareClock, err := strconv.Atoi(strings.TrimSpace(output))
	if err != nil {
		return -1, fmt.Errorf("failed to convert PTP hardware clock for interface %s on node %s to int: %w",
			ifName, nodeName, err)
	}

	return hardwareClock, nil
}

// AdjustPTPHardwareClock adjusts the PTP hardware clock for a given network interface on a specified node. This affects
// the CLOCK_REALTIME offset. The amount is in seconds.
func AdjustPTPHardwareClock(client *clients.Settings, nodeName string, ifName Name, amount float64) error {
	hardwareClock, err := GetPTPHardwareClock(client, nodeName, ifName)
	if err != nil {
		return fmt.Errorf("failed to get PTP hardware clock for interface %s on node %s: %w", ifName, nodeName, err)
	}

	command := fmt.Sprintf("phc_ctl /dev/ptp%d adjust %f", hardwareClock, amount)
	_, err = ptpdaemon.ExecuteCommandInPtpDaemonPod(client, nodeName, command)

	if err != nil {
		return fmt.Errorf("failed to adjust PTP hardware clock for interface %s on node %s: %w", ifName, nodeName, err)
	}

	return nil
}

// ResetPTPHardwareClock resets the PTP hardware clock for a given network interface on a specified node.
func ResetPTPHardwareClock(client *clients.Settings, nodeName string, ifName Name) error {
	hardwareClock, err := GetPTPHardwareClock(client, nodeName, ifName)
	if err != nil {
		return fmt.Errorf("failed to get PTP hardware clock for interface %s on node %s: %w", ifName, nodeName, err)
	}

	command := fmt.Sprintf("phc_ctl /dev/ptp%d set", hardwareClock)
	_, err = ptpdaemon.ExecuteCommandInPtpDaemonPod(client, nodeName, command)

	if err != nil {
		return fmt.Errorf("failed to reset PTP hardware clock for interface %s on node %s: %w", ifName, nodeName, err)
	}

	return nil
}
