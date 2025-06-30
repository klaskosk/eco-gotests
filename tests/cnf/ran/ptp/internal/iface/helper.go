package iface

import (
	"fmt"
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
