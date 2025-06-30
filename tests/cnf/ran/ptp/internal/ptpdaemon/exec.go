// Package ptpdaemon provides functions for executing commands in the PTP daemon pod.
package ptpdaemon

import (
	"fmt"

	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// GetPtpDaemonPodOnNode retrieves the PTP daemon pod running on the specified node. It returns an error if it cannot
// find exactly one PTP daemon pod on the node.
func GetPtpDaemonPodOnNode(client *clients.Settings, nodeName string) (*pod.Builder, error) {
	daemonPods, err := pod.List(client, ranparam.PtpOperatorNamespace, metav1.ListOptions{
		LabelSelector: ranparam.PtpDaemonsetLabelSelector,
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": nodeName}).String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list PTP daemon pods on node %s: %w", nodeName, err)
	}

	if len(daemonPods) != 1 {
		return nil, fmt.Errorf("expected exactly one PTP daemon pod on node %s, found %d", nodeName, len(daemonPods))
	}

	return daemonPods[0], nil
}

// ExecuteCommandInPtpDaemonPod executes a command in the PTP daemon pod running on the specified node. It returns the
// output of the command as a string. If the command execution fails, it returns an error.
func ExecuteCommandInPtpDaemonPod(client *clients.Settings, nodeName string, command string) (string, error) {
	daemonPod, err := GetPtpDaemonPodOnNode(client, nodeName)
	if err != nil {
		return "", fmt.Errorf("failed to get PTP daemon pod on node %s: %w", nodeName, err)
	}

	output, err := daemonPod.ExecCommand([]string{"sh", "-c", command}, ranparam.PtpContainerName)
	if err != nil {
		return "", fmt.Errorf("failed to execute command `%v` in PTP daemon pod on node %s: %w", command, nodeName, err)
	}

	return output.String(), nil
}
