package ranhelper

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/deployment"
	"github.com/openshift-kni/eco-goinfra/pkg/nto"
	"github.com/openshift-kni/eco-goinfra/pkg/olm"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
)

// InitializeSpokeNames initializes the name of spoke 1 and, if present, spoke 2.
func InitializeSpokeNames() error {
	var err error

	// Spoke 1 is required to be present.
	ranparam.Spoke1Name, err = getClusterName(os.Getenv("KUBECONFIG"))
	if err != nil {
		return err
	}

	// Spoke 2 is optional depending on the test.
	if raninittools.RANConfig.Spoke2Kubeconfig != "" {
		ranparam.Spoke2Name, err = getClusterName(raninittools.RANConfig.Spoke2Kubeconfig)
		if err != nil {
			return err
		}
	}

	return nil
}

// InitializeTalmVersion initializes the version of the TALM operator from the hub cluster.
func InitializeTalmVersion() error {
	var err error

	ranparam.TalmVersion, err = getOperatorVersionFromCsv(
		raninittools.HubAPIClient, ranparam.TalmOperatorHubNamespace, ranparam.OpenshiftOperatorNamespace)

	return err
}

// InitializeVersions initializes the versions of ACM, TALM, and ZTP from the hub cluster.
func InitializeVersions() error {
	var err error

	err = InitializeTalmVersion()
	if err != nil {
		return err
	}

	ranparam.AcmVersion, err = getOperatorVersionFromCsv(
		raninittools.HubAPIClient, ranparam.AcmOperatorName, ranparam.AcmOperatorNamespace)
	if err != nil {
		return err
	}

	ranparam.ZtpVersion, err = getZtpVersionFromArgoCd(
		raninittools.HubAPIClient, ranparam.OpenshiftGitopsRepoServer, ranparam.OpenshiftGitops)

	return err
}

// IsPodHealthy returns true if a given pod is healthy, otherwise false.
func IsPodHealthy(pod *pod.Builder) bool {
	if pod.Object.Status.Phase == corev1.PodRunning {
		// Check if running pod is ready
		if !isPodInCondition(pod, corev1.PodReady) {
			glog.V(ranparam.LogLevel).Infof("pod condition is not Ready. Message: %s", pod.Object.Status.Message)

			return false
		}
	} else if pod.Object.Status.Phase != corev1.PodSucceeded {
		// Pod is not running or completed.
		glog.V(ranparam.LogLevel).Infof("pod phase is %s. Message: %s", pod.Object.Status.Phase, pod.Object.Status.Message)

		return false
	}

	return true
}

// DoesContainerExistInPod checks if a given container exists in a given pod.
func DoesContainerExistInPod(pod *pod.Builder, containerName string) bool {
	containers := pod.Object.Status.ContainerStatuses

	for _, container := range containers {
		if container.Name == containerName {
			glog.V(ranparam.LogLevel).Infof("found %s container", containerName)

			return true
		}
	}

	return false
}

// IsVersionStringInRange checks if a version string is between a specified min and max value, inclusive. All the string
// inputs to this function should be dot separated positive intergers, e.g. "1.0.0" or "4.10". Each string input must be
// at least two dot separarted integers but may also be 3 or more, though only the first two are compared.
func IsVersionStringInRange(version, minimum, maximum string) (bool, error) {
	versionValid, versionDigits := validateInputString(version)
	minimumValid, minimumDigits := validateInputString(minimum)
	maximumValid, maximumDigits := validateInputString(maximum)

	if !minimumValid {
		// Only accept invalid empty strings
		if minimum != "" {
			return false, fmt.Errorf("invalid minimum provided: '%s'", minimum)
		}

		// Assume the minimum digits are [0,0] for later comparison
		minimumDigits = []int{0, 0}
	}

	if !maximumValid {
		// Only accept invalid empty strings
		if maximum != "" {
			return false, fmt.Errorf("invalid maximum provided: '%s'", maximum)
		}

		// Assume the maximum digits are [math.MaxInt, math.MaxInt] for later comparison
		maximumDigits = []int{math.MaxInt, math.MaxInt}
	}

	// If the version was not valid then we need to check the min and max
	if !versionValid {
		// If no min or max was defined then return true
		if !minimumValid && !maximumValid {
			return true, nil
		}

		// Otherwise return whether the input maximum was an empty string or not
		return maximum == "", nil
	}

	// Otherwise the versions were valid so compare the digits
	for i := 0; i < 2; i++ {
		// The version bit should be between the minimum and maximum
		if versionDigits[i] < minimumDigits[i] || versionDigits[i] > maximumDigits[i] {
			return false, nil
		}
	}

	// At the end if we never returned then all the digits were in valid range
	return true, nil
}

type SelectorFunc func(profile *nto.Builder) bool

// GetPerformanceProfileWithCPUSet returns the first performance profile found with reserved and isolated cpuset. If a
// selector function is provided, this becomes an additional condition that the performance profile must fulfil.
func GetPerformanceProfileWithCPUSet(client *clients.Settings, selectorFunc ...SelectorFunc) (*nto.Builder, error) {
	profileBuilders, err := nto.ListProfiles(client)
	if err != nil {
		return nil, err
	}

	for _, profileBuilder := range profileBuilders {
		if profileBuilder.Object.Spec.CPU != nil &&
			profileBuilder.Object.Spec.CPU.Reserved != nil &&
			profileBuilder.Object.Spec.CPU.Isolated != nil {
			if len(selectorFunc) > 0 && !selectorFunc[0](profileBuilder) {
				continue
			}

			return profileBuilder, nil
		}
	}

	return nil, errors.New("failed to find performance profile with reserved and isolated CPU set")
}

// UnmarshalRaw converts raw bytes for a K8s CR into the actual type.
func UnmarshalRaw[T any](raw []byte) (*T, error) {
	untyped := &unstructured.Unstructured{}
	err := untyped.UnmarshalJSON(raw)

	if err != nil {
		return nil, err
	}

	var typed T
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(untyped.UnstructuredContent(), &typed)

	if err != nil {
		return nil, err
	}

	return &typed, nil
}

// AreClustersPresent checks all of the provided clusters and returns false if any are nil.
func AreClustersPresent(clusters []*clients.Settings) bool {
	for _, cluster := range clusters {
		if cluster == nil {
			return false
		}
	}

	return true
}

// validateInputString validates that a string is at least two dot separated nonnegative integers.
func validateInputString(input string) (bool, []int) {
	versionSplits := strings.Split(input, ".")

	if len(versionSplits) < 2 {
		return false, []int{}
	}

	digits := []int{}

	for i := 0; i < 2; i++ {
		digit, err := strconv.Atoi(versionSplits[i])
		if err != nil || digit < 0 {
			return false, []int{}
		}

		digits = append(digits, digit)
	}

	return true, digits
}

// isPodInCondition returns true if a given pod is in expected condition, otherwise false.
func isPodInCondition(pod *pod.Builder, condition corev1.PodConditionType) bool {
	for _, c := range pod.Object.Status.Conditions {
		if c.Type == condition && c.Status == corev1.ConditionTrue {
			return true
		}
	}

	return false
}

// getClusterName extracts the cluster name from provided kubeconfig, assuming there's one cluster in the kubeconfig.
func getClusterName(kubeconfigPath string) (string, error) {
	rawConfig, _ := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()

	for _, cluster := range rawConfig.Clusters {
		// Get a cluster name by parsing it from the server hostname. Expects the url to start with
		// `https://api.cluster-name.` so splitting by `.` gives the cluster name.
		splits := strings.Split(cluster.Server, ".")
		clusterName := splits[1]

		glog.V(ranparam.LogLevel).Infof("cluster name %s found for kubeconfig at %s", clusterName, kubeconfigPath)

		return clusterName, nil
	}

	return "", fmt.Errorf("could not get cluster name for kubeconfig at %s", kubeconfigPath)
}

// getOperatorVersionFromCsv returns operator version from csv, or an empty string if no CSV for the provided operator
// is found.
func getOperatorVersionFromCsv(client *clients.Settings, operatorName, operatorNamespace string) (string, error) {
	csv, err := olm.ListClusterServiceVersion(client, operatorNamespace, metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	for _, csv := range csv {
		if strings.Contains(csv.Object.Name, operatorName) {
			return csv.Object.Spec.Version.String(), nil
		}
	}

	return "", fmt.Errorf("could not find version for operator %s in namespace %s", operatorName, operatorNamespace)
}

// getZtpVersionFromArgoCd is used to fetch the version of the ztp-site-generate init container.
func getZtpVersionFromArgoCd(client *clients.Settings, name, namespace string) (string, error) {
	ztpDeployment, err := deployment.Pull(client, name, namespace)
	if err != nil {
		return "", err
	}

	for _, container := range ztpDeployment.Definition.Spec.Template.Spec.InitContainers {
		// Match both the `ztp-site-generator` and `ztp-site-generate` images since which one depends on
		// versions.
		if strings.Contains(container.Image, "ztp-site-gen") {
			colonSplit := strings.Split(container.Image, ":")
			ztpVersion := colonSplit[len(colonSplit)-1]

			if ztpVersion == "latest" {
				glog.V(ranparam.LogLevel).Info("ztp-site-generate version tag was 'latest', returning empty version")

				return "", nil
			}

			// The format here will be like vX.Y.Z so we need to remove the v at the start.
			return ztpVersion[1:], nil
		}
	}

	return "", errors.New("unable to identify ZTP version")
}
