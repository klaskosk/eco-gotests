package helper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/mco"
	"github.com/openshift-kni/eco-goinfra/pkg/nodes"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/cluster"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/tsparams"
	mcv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/cpuset"
)

// containerInfoCommand is the command to get information about running containers.
const containerInfoCommand = `crictl ps --state running --quiet | xargs crictl inspect -o json | jq '. | {
name: .status.metadata.name,
podname: .status.labels."io.kubernetes.pod.name",
namespace: .status.labels."io.kubernetes.pod.namespace",
pid: .info.pid,
cpus: .info.runtimeSpec.linux.resources.cpu.cpus,
shares: .info.runtimeSpec.linux.resources.cpu.shares,
}' | jq -sM`

// pidAndAffinityRegexp is a regex to parse the PID and affinity list from the output of the taskset command.
var pidAndAffinityRegexp = regexp.MustCompile(`pid (\d+)'s current affinity list: (.*)$`)

// WaitForNodeFunctional waits for the MCP to be updating, then updated, then all nodes ready.
func WaitForNodeFunctional(client *clients.Settings) error {
	mcp, err := mco.Pull(raninittools.Spoke1APIClient, tsparams.MCPName)
	if err != nil {
		return err
	}

	err = mcp.WaitToBeInCondition(mcv1.MachineConfigPoolUpdating, corev1.ConditionTrue, 10*time.Minute)
	if err != nil {
		return err
	}

	err = mcp.WaitForUpdate(time.Hour)
	if err != nil {
		return err
	}

	nodeList, err := nodes.List(client)
	if err != nil {
		return err
	}

	for _, node := range nodeList {
		err := node.WaitUntilReady(10 * time.Minute)
		if err != nil {
			return err
		}
	}

	return nil
}

// ContainerInfo contains information about a running container.
type ContainerInfo struct {
	Name      string `json:"name"`
	Cpus      string `json:"cpus"`
	Namespace string `json:"namespace"`
	PodName   string `json:"podname"`
	Shares    int    `json:"shares"`
	Pid       int    `json:"pid"`
}

// CheckAffinitiesByProcessMatch checks CPU affinities for an array of processes against an specfied reserved CPU set.
func CheckAffinitiesByProcessMatch(
	client *clients.Settings, processNames []string, reservedCPUSet cpuset.CPUSet) error {
	for _, processName := range processNames {
		cmd := fmt.Sprintf("pgrep %s | while read i; do taskset -cp $i; done", processName)
		output, err := cluster.ExecCommandOnSNO(client, 3, cmd)

		if err != nil {
			return err
		}

		for _, line := range strings.Split(output, "\r\n") {
			line = strings.TrimSpace(line)

			if len(line) == 0 {
				return fmt.Errorf("process name: %s is not matched", processName)
			}

			match := pidAndAffinityRegexp.FindAllStringSubmatch(line, -1)
			if match == nil {
				return fmt.Errorf("unmatched pid and affinity for process name: %s", processName)
			}

			pid := match[0][1]
			affinity := strings.TrimSpace(match[0][2])
			pidCpuset, err := cpuset.Parse(affinity)

			if err != nil {
				return err
			}

			if !pidCpuset.IsSubsetOf(reservedCPUSet) {
				return fmt.Errorf("process: %s pid: %s with actual affinity: %s but expected: %s",
					processName, pid, affinity, reservedCPUSet.String())
			}
		}
	}

	return nil
}

// GetContainersInfo returns containers info on given node via crictl, trying up to 3 times.
func GetContainersInfo(client *clients.Settings) ([]ContainerInfo, error) {
	var (
		containerInfos []ContainerInfo
		output         string
		err            error
	)

	// Retry up to 3 times in case the output is formatted wrong.
	for range 3 {
		output, err = cluster.ExecCommandOnSNO(client, 3, containerInfoCommand)
		if err == nil {
			err = json.Unmarshal([]byte(output), &containerInfos)
		}

		if err == nil {
			break
		}

		glog.V(tsparams.LogLevel).Info("Error getting container info, retrying")

		time.Sleep(10 * time.Second)
	}

	return containerInfos, err
}

// GetManagementContainersInfo filters out test pods from the provided container infos.
func GetManagementContainersInfo(containerInfos []ContainerInfo) []ContainerInfo {
	var managementInfos []ContainerInfo

	for _, containerInfo := range containerInfos {
		if !tsparams.NonManagementNamespaces.Has(containerInfo.Namespace) &&
			!strings.HasPrefix(containerInfo.PodName, "cnf-ran-gotests-priv") &&
			!strings.HasPrefix(containerInfo.PodName, "process-exporter") {
			managementInfos = append(managementInfos, containerInfo)
		}
	}

	return managementInfos
}

// CheckPodsAffinity checks given containers are pinned to specified cpus.
func CheckPodsAffinity(containersInfo []ContainerInfo, affinedCPUSet cpuset.CPUSet) error {
	for _, containerInfo := range containersInfo {
		cpus, err := cpuset.Parse(containerInfo.Cpus)
		if err != nil {
			return err
		}

		if !cpus.Equals(affinedCPUSet) {
			return fmt.Errorf("found container not pinned to specified cpus")
		}
	}

	return nil
}

// CheckCPUAffinityOnNonKernelPids checks the CPU affinity on all non-kernel PIDs and if any are not confined to the
// specified CPUSet an error is returned along with the map of the PIDs not confined to the CPUSet and their affinities.
func CheckCPUAffinityOnNonKernelPids(client *clients.Settings, cpus cpuset.CPUSet) (map[int]string, error) {
	pidsToExclude, err := getKernelPids(client)
	if err != nil {
		return nil, err
	}

	fecPids, err := getFecPids(client)
	if err != nil {
		return nil, err
	}

	pidsToExclude = append(pidsToExclude, fecPids...)

	allPids, err := getAllPids(client)
	if err != nil {
		return nil, err
	}

	containersInfo, err := GetContainersInfo(client)
	if err != nil {
		return nil, err
	}

	for _, containerInfo := range containersInfo {
		pidsToExclude = append(pidsToExclude, containerInfo.Pid)
	}

	var pidsToCheck []int

	for _, pid := range allPids {
		if !slices.Contains(pidsToExclude, pid) {
			pidsToCheck = append(pidsToCheck, pid)
		}
	}

	affinities, err := getPidsAffinity(client, pidsToCheck)
	if err != nil {
		return nil, err
	}

	failedMap := make(map[int]string)

	var failedPids []int

	for pid, affinity := range affinities {
		pidCpuset, err := cpuset.Parse(affinity)
		if err != nil {
			return nil, err
		}

		if !pidCpuset.IsSubsetOf(cpus) {
			// Make sure it's not a kernel process as there may be a race condition in the previous queries
			// of pids.
			if !isKernelPid(client, pid) {
				failedMap[pid] = affinity
				failedPids = append(failedPids, pid)
			}
		}
	}

	if len(failedPids) > 0 {
		glog.V(tsparams.LogLevel).Info(getPidInfo(client, failedPids))

		err = fmt.Errorf("processess not matching reserved CPU affinities found")
	}

	return failedMap, err
}

// getPidInfo uses the ps command to return information about all of the provided pids.
func getPidInfo(client *clients.Settings, pids []int) string {
	var pidStrings []string
	for _, pid := range pids {
		pidStrings = append(pidStrings, strconv.Itoa(pid))
	}

	cmd := fmt.Sprintf("ps %s", strings.Join(pidStrings, " "))
	output, _ := cluster.ExecCommandOnSNO(client, 3, cmd)

	return output
}

// isKernelPid checks if the given process is a kernel process and returns true if it is or if there was an error
// getting the pid, likely due to the process being terminated.
func isKernelPid(client *clients.Settings, pid int) bool {
	cmd := fmt.Sprintf("ps --no-headers -o ppid -p %d > /tmp/x ; cat /tmp/x", pid)
	parentPid, _ := cluster.ExecCommandOnSNO(client, 3, cmd)

	if len(parentPid) == 0 || strings.TrimSpace(parentPid) == "2" {
		// Either the process was terminated or it is a kernel process.
		return true
	}

	return false
}

// getPids runs the provided command to get a list of pids, trying up to 3 times when encountering Signal 23 errors.
func getPids(client *clients.Settings, command string) ([]int, error) {
	var (
		output string
		err    error
	)

	// ps command via container with long output often causes SIGURG, so redirect output to a file first.
	command += " > /tmp/x ; cat /tmp/x"

	// Retry in case of SIGURG output.
	for range 3 {
		output, err = cluster.ExecCommandOnSNO(client, 3, command)
		if err != nil {
			return nil, err
		}

		if !strings.Contains(output, "Signal 23") {
			break
		}

		glog.V(tsparams.LogLevel).Info("Signal 23 encountered, retrying")

		time.Sleep(10 * time.Second)
	}

	var pids []int

	if len(strings.TrimSpace(output)) == 0 {
		return pids, nil
	}

	for _, pidString := range strings.Split(output, "\r\n") {
		pidString = strings.TrimSpace(pidString)
		pid, err := strconv.Atoi(pidString)

		if err != nil {
			return nil, err
		}

		pids = append(pids, pid)
	}

	return pids, nil
}

// getKernelPids returns list of kernel process ids.
func getKernelPids(client *clients.Settings) ([]int, error) {
	// Get all kernel threads (PID 2 and children)
	return getPids(client, "ps --no-headers --ppid 2 -p 2 -o pid")
}

// getFecPids returns list of "/sriov_workdir/pf_bb_config" process ids used by FEC.
func getFecPids(client *clients.Settings) ([]int, error) {
	return getPids(client, "pgrep pf_bb_config")
}

// getAllPids returns list of all process ids.
func getAllPids(client *clients.Settings) ([]int, error) {
	return getPids(client, "ps --no-headers -e -o pid")
}

// getPidsAffinity gets pids' affinity list from taskset command and returns a map with pid as key and affinity list as
// value.
func getPidsAffinity(client *clients.Settings, pids []int) (map[int]string, error) {
	var pidStrings []string

	for _, pid := range pids {
		pidStrings = append(pidStrings, strconv.Itoa(pid))
	}

	pidString := strings.Join(pidStrings, " ")
	cmd := fmt.Sprintf(
		"pids=\"%s\"; for pid in $pids; do if [ $pid == $$ ]; then continue; fi; taskset -pc $pid; done", pidString)
	output, err := cluster.ExecCommandOnSNO(client, 3, cmd)

	if err != nil {
		return nil, err
	}

	if !strings.Contains(output, "current affinity list") {
		return nil, fmt.Errorf("current affinity list not found in command output")
	}

	affinities := make(map[int]string)

	for _, line := range strings.Split(output, "\r\n") {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "No such process") {
			continue
		}

		match := pidAndAffinityRegexp.FindAllStringSubmatch(line, -1)
		pidInt, err := strconv.Atoi(match[0][1])

		if err != nil {
			return nil, err
		}

		affinities[pidInt] = strings.TrimSpace(match[0][2])
	}

	return affinities, nil
}
