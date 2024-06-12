package helper

import (
	"errors"
	"fmt"

	"github.com/openshift-kni/eco-goinfra/pkg/namespace"
	"github.com/openshift-kni/eco-goinfra/pkg/pod"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranhelper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/talm/internal/tsparams"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateTalmTestNamespace creates the test namespace on the hub and the spokes, if present.
func CreateTalmTestNamespace() error {
	_, err := namespace.NewBuilder(raninittools.Spoke1APIClient, tsparams.TestNamespace).Create()
	if err != nil {
		return err
	}

	if raninittools.HubAPIClient != nil {
		_, err := namespace.NewBuilder(raninittools.HubAPIClient, tsparams.TestNamespace).Create()
		if err != nil {
			return err
		}
	}

	// spoke 2 may be optional depending on what tests are running
	if raninittools.Spoke2APIClient != nil {
		_, err := namespace.NewBuilder(raninittools.Spoke2APIClient, tsparams.TestNamespace).Create()
		if err != nil {
			return err
		}
	}

	return nil
}

// VerifyTalmIsInstalled checks that talm pod and container is present and that CGUs can be fetched.
func VerifyTalmIsInstalled() error {
	// Check for talm pods
	talmPods, err := pod.List(
		raninittools.HubAPIClient,
		ranparam.OpenshiftOperatorNamespace,
		metav1.ListOptions{LabelSelector: tsparams.TalmPodLabelSelector})
	if err != nil {
		return err
	}

	// Check if any pods exist
	if len(talmPods) == 0 {
		return errors.New("unable to find talm pod")
	}

	// Check each pod for the talm container
	for _, talmPod := range talmPods {
		if !ranhelper.IsPodHealthy(talmPod) {
			return fmt.Errorf("talm pod %s is not healthy", talmPod.Definition.Name)
		}

		if !ranhelper.DoesContainerExistInPod(talmPod, ranparam.TalmContainerName) {
			return errors.New("talm pod defined but talm container does not exist")
		}
	}

	return nil
}
