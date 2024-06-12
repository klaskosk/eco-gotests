package helper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/argocd"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/imageregistry"
	"github.com/openshift-kni/eco-goinfra/pkg/ocm"
	"github.com/openshift-kni/eco-goinfra/pkg/serviceaccount"
	"github.com/openshift-kni/eco-goinfra/pkg/storage"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranhelper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/ztp/internal/tsparams"
	operatorv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	configurationPolicyV1 "open-cluster-management.io/config-policy-controller/api/v1"
)

const (
	// LabelRole contains the key for the role label.
	LabelRole = "node-role.kubernetes.io"
)

// WaitForServiceAccountToExist waits for up to the specified timeout until the service account exists.
func WaitForServiceAccountToExist(client *clients.Settings, name, namespace string, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(
		context.TODO(), tsparams.ArgoCdChangeInterval, timeout, true, func(ctx context.Context) (bool, error) {
			_, err := serviceaccount.Pull(client, name, namespace)
			if err == nil {
				return true, nil
			}

			if strings.Contains(err.Error(), "does not exist") {
				return false, nil
			}

			return false, err
		})
}

// GetPolicyEvaluationIntervals is used to get the configured evaluation intervals for the specified policy.
func GetPolicyEvaluationIntervals(policy *ocm.PolicyBuilder) (string, string, error) {
	glog.V(tsparams.LogLevel).Infof(
		"Checking policy '%s' in namespace '%s' to fetch evaluation intervals",
		policy.Definition.Name, policy.Definition.Namespace)

	policyTemplates := policy.Definition.Spec.PolicyTemplates
	if len(policyTemplates) < 1 {
		return "", "", fmt.Errorf(
			"could not find policy template for policy %s/%s", policy.Definition.Name, policy.Definition.Namespace)
	}

	configPolicy, err := ranhelper.UnmarshalRaw[configurationPolicyV1.ConfigurationPolicy](policyTemplates[0].ObjectDefinition.Raw)
	if err != nil {
		return "", "", err
	}

	complianceInterval := configPolicy.Spec.EvaluationInterval.Compliant
	nonComplianceInterval := configPolicy.Spec.EvaluationInterval.NonCompliant

	return complianceInterval, nonComplianceInterval, nil
}

// WaitForConditionInArgoCdApp waits up to timeout until the specified Argo CD app has a condition containing the
// expectedMessage.
func WaitForConditionInArgoCdApp(
	client *clients.Settings, appName, namespace, expectedMessage string, timeout time.Duration) error {
	glog.V(tsparams.LogLevel).Infof(
		"Checking application '%s' in namespace '%s' for condition with message '%s'", appName, namespace, expectedMessage)

	return wait.PollUntilContextTimeout(
		context.TODO(), tsparams.ArgoCdChangeInterval, timeout, true, func(ctx context.Context) (bool, error) {
			app, err := argocd.PullApplication(client, appName, namespace)
			if err != nil {
				return false, err
			}

			for _, condition := range app.Definition.Status.Conditions {
				if strings.Contains(condition.Message, expectedMessage) {
					glog.V(tsparams.LogLevel).Info("Found matching condition")

					return true, nil
				}

				glog.V(tsparams.LogLevel).Infof("Condition message '%s' did not match", condition.Message)
			}

			return false, nil
		})
}

// RestoreImageRegistry restores the image registry with the provided name back to imageRegistryConfig, copying over the
// labels, annotations, and spec from imageRegistryConfig, then waiting until the image registry is available again.
func RestoreImageRegistry(
	client *clients.Settings, imageRegistryName string, imageRegistryConfig *imageregistry.Builder) error {
	currentImageRegistry, err := imageregistry.Pull(client, imageRegistryName)
	if err != nil {
		return err
	}

	if imageRegistryConfig.Definition.GetAnnotations() != nil {
		currentImageRegistry.Definition.SetAnnotations(imageRegistryConfig.Definition.GetAnnotations())
	}

	if imageRegistryConfig.Definition.GetLabels() != nil {
		currentImageRegistry.Definition.SetLabels(imageRegistryConfig.Definition.GetLabels())
	}

	currentImageRegistry.Definition.Spec = imageRegistryConfig.Definition.Spec

	currentImageRegistry, err = currentImageRegistry.Update()
	if err != nil {
		return err
	}

	return WaitForConditionInImageRegistry(currentImageRegistry, metav1.Condition{
		Type:   "Available",
		Reason: "Removed",
		Status: metav1.ConditionTrue,
	}, tsparams.ArgoCdChangeTimeout)
}

// WaitForConditionInImageRegistry waits until the image registry has a condition that matches the expected,
// checking only the Type, Status, Reason, and Message fields. Zero fields in the expected condition are ignored.
func WaitForConditionInImageRegistry(
	imageRegistryConfig *imageregistry.Builder, expected metav1.Condition, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(
		context.TODO(), tsparams.ArgoCdChangeInterval, timeout, true, func(ctx context.Context) (bool, error) {
			if !imageRegistryConfig.Exists() {
				glog.V(tsparams.LogLevel).Infof("imageRegistry %s does not exist in namespace %s",
					imageRegistryConfig.Definition.Name, imageRegistryConfig.Definition.Namespace)

				return false, nil
			}

			for _, condition := range imageRegistryConfig.Definition.Status.Conditions {
				if (expected.Type != "" && condition.Type != expected.Type) ||
					(expected.Status != "" && condition.Status != operatorv1.ConditionStatus(expected.Status)) ||
					(expected.Reason != "" && condition.Reason != expected.Reason) ||
					(expected.Message != "" && condition.Message != expected.Message) {
					continue
				}

				return true, nil
			}

			return false, nil
		})
}

// CleanupImageRegistryConfig deletes the specified resources in the necessary order.
func CleanupImageRegistryConfig(
	client *clients.Settings,
	storageClassName,
	persistentVolumeName,
	persistentVolumeClaimName,
	persistentVolumeClaimNamespace string) error {
	glog.V(tsparams.LogLevel).Infof(
		"Cleaning up image registry resources with sc=%s, pv=%s, pvc=%s",
		storageClassName, persistentVolumeName, persistentVolumeClaimName)

	// The resources must be deleted in the order of pvc, pv, then sc to avoid errors.
	if persistentVolumeClaimName != "" {
		pvc, err := storage.PullPersistentVolumeClaim(client, persistentVolumeClaimName, persistentVolumeClaimNamespace)
		if err != nil {
			return err
		}

		err = pvc.DeleteAndWait(tsparams.ArgoCdChangeTimeout)
		if err != nil {
			return err
		}
	}

	if persistentVolumeName != "" {
		_, err := storage.PullPersistentVolume(client, persistentVolumeName)
		if err != nil {
			return err
		}

		return fmt.Errorf("waiting on eco-goinfra PR for pv delete")
	}

	if storageClassName != "" {
		return fmt.Errorf("need to write eco-goinfra pr for storageclass pull")
	}

	return nil
}
