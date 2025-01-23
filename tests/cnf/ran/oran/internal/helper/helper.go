package helper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/bmc"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/ocm"
	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	"github.com/openshift-kni/eco-goinfra/pkg/siteconfig"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	pluginv1alpha1 "github.com/openshift-kni/oran-hwmgr-plugin/api/hwmgr-plugin/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	policiesv1 "open-cluster-management.io/governance-policy-propagator/api/v1"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// NewProvisioningRequest creates a ProvisioningRequest builder using the provided values. It sets the
// templateParameters automatically.
func NewProvisioningRequest(
	client *clients.Settings, clusterName, hostName, templateVersion string) *oran.ProvisioningRequestBuilder {
	prBuilder := oran.NewPRBuilder(client, tsparams.TestPRName, tsparams.ClusterTemplateName, templateVersion).
		WithTemplateParameter("nodeClusterName", clusterName).
		WithTemplateParameter("oCloudSiteId", clusterName).
		WithTemplateParameter("policyTemplateParameters", map[string]any{}).
		WithTemplateParameter("clusterInstanceParameters", map[string]any{
			"clusterName": clusterName,
			"nodes": []map[string]any{{
				"hostName": hostName,
			}},
		})

	return prBuilder
}

// NewNoTemplatePR creates a ProvisioningRequest builder using the provided values, following the schema for no
// HardwareTemplate. The BMC and network data is incorrect so that a ClusterInstance is generated but will not actually
// provision.
func NewNoTemplatePR(
	client *clients.Settings, clusterName, hostName, templateVersion string) *oran.ProvisioningRequestBuilder {
	prBuilder := oran.NewPRBuilder(client, tsparams.TestPRName, tsparams.ClusterTemplateName, templateVersion).
		WithTemplateParameter("nodeClusterName", clusterName).
		WithTemplateParameter("oCloudSiteId", clusterName).
		WithTemplateParameter("policyTemplateParameters", map[string]any{}).
		WithTemplateParameter("clusterInstanceParameters", map[string]any{
			"clusterName": clusterName,
			"nodes": []map[string]any{{
				"hostName":   hostName,
				"bmcAddress": "redfish-VirtualMedia://10.10.10.10/redfish/v1/Systems/System.Embedded.1",
				"bmcCredentialsDetails": map[string]any{
					"username": tsparams.TestBase64Credential,
					"password": tsparams.TestBase64Credential,
				},
				"bootMACAddress": "01:23:45:67:89:AB",
				"nodeNetwork": map[string]any{
					"interfaces": []map[string]any{{
						"macAddress": "01:23:45:67:89:AB",
					}},
				},
			}},
		})

	return prBuilder
}

// GetValidDellHwmgr returns the first HardwareManager with AdaptorID dell-hwmgr and where condition Validation is True.
func GetValidDellHwmgr(client *clients.Settings) (*oran.HardwareManagerBuilder, error) {
	hwmgrs, err := oran.ListHardwareManagers(client, runtimeclient.ListOptions{
		Namespace: tsparams.HardwareManagerNamespace,
	})
	if err != nil {
		return nil, err
	}

	for _, hwmgr := range hwmgrs {
		if hwmgr.Definition.Spec.AdaptorID != pluginv1alpha1.SupportedAdaptors.Dell {
			continue
		}

		for _, condition := range hwmgr.Definition.Status.Conditions {
			if condition.Type == string(pluginv1alpha1.ConditionTypes.Validation) && condition.Status == metav1.ConditionTrue {
				return hwmgr, nil
			}
		}
	}

	return nil, fmt.Errorf("no valid HardwareManager with AdaptorID dell-hwmgr exists")
}

// WaitForCIExtraLabel waits up to timeout until clusterInstance contains the label. This will update the Object but not
// the Definition of the CIBuilder.
func WaitForCIExtraLabel(clusterInstance *siteconfig.CIBuilder, label string, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(
		context.TODO(), 3*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
			var err error
			clusterInstance.Object, err = clusterInstance.Get()

			if err != nil {
				glog.V(tsparams.LogLevel).Infof("Failed to get ClusterInstance %s in namespace %s: %v",
					clusterInstance.Definition.Name, clusterInstance.Definition.Namespace)

				return false, nil
			}

			mclLabels, ok := clusterInstance.Object.Spec.ExtraLabels["ManagedCluster"]
			if !ok {
				glog.V(tsparams.LogLevel).Infof("ClusterInstance %s in namespace %s is missing ManagedCluster extraLabels",
					clusterInstance.Definition.Name, clusterInstance.Definition.Namespace)

				return false, nil
			}

			_, containsLabel := mclLabels[label]

			return containsLabel, nil
		})
}

// WaitForMCLLabel waits up to timeout until mcl contains the label. This will update the Object but not the Definition
// of the ManagedClusterBuilder.
func WaitForMCLLabel(mcl *ocm.ManagedClusterBuilder, label string, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(
		context.TODO(), 3*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
			if !mcl.Exists() {
				glog.V(tsparams.LogLevel).Infof("ManagedCluster %s does not exist", mcl.Definition.Name)

				return false, nil
			}

			if mcl.Object == nil {
				glog.V(tsparams.LogLevel).Infof("Failed to get ManagedCluster %s", mcl.Definition.Name)

				return false, nil
			}

			_, containsLabel := mcl.Object.Labels[label]

			return containsLabel, nil
		})
}

// WaitForPoliciesCompliant waits up to the timeout until all of the policies in namespace are Compliant.
func WaitForPoliciesCompliant(client *clients.Settings, namespace string, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(
		context.TODO(), 3*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
			policies, err := ocm.ListPoliciesInAllNamespaces(client, runtimeclient.ListOptions{Namespace: namespace})
			if err != nil {
				glog.V(tsparams.LogLevel).Infof("Failed to list all policies in namespace %s: %v", namespace, err)

				return false, nil
			}

			for _, policy := range policies {
				if policy.Definition.Status.ComplianceState != policiesv1.Compliant {
					glog.V(tsparams.LogLevel).Infof("Policy %s in namespace %s is not compliant",
						policy.Definition.Name, policy.Definition.Namespace)

					return false, nil
				}
			}

			return true, nil
		})
}

// WaitForNoncompliantImmutable waits up to timeout until one of the policies in namespace is NonCompliant and the
// message history shows it is due to an immutable field.
func WaitForNoncompliantImmutable(client *clients.Settings, namespace string, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(
		context.TODO(), 3*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
			policies, err := ocm.ListPoliciesInAllNamespaces(client, runtimeclient.ListOptions{Namespace: namespace})
			if err != nil {
				glog.V(tsparams.LogLevel).Infof("Failed to list all policies in namespace %s: %v", namespace, err)

				return false, nil
			}

			for _, policy := range policies {
				if policy.Definition.Status.ComplianceState == policiesv1.NonCompliant {
					glog.V(tsparams.LogLevel).Infof("Policy %s in namespace %s is not compliant, checking history",
						policy.Definition.Name, policy.Definition.Namespace)

					details := policy.Definition.Status.Details
					if len(details) != 1 {
						continue
					}

					history := details[0].History
					if len(history) < 1 {
						continue
					}

					if strings.Contains(history[0].Message, tsparams.ImmutableMessage) {
						glog.V(tsparams.LogLevel).Infof("Policy %s in namespace %s is not compliant due to an immutable field",
							policy.Definition.Name, policy.Definition.Namespace)

						return true, nil
					}
				}
			}

			return false, nil
		})
}

// WaitForPoweredOff waits up to timeout until the provided BMC shows the system is powered off.
func WaitForPoweredOff(bmcClient *bmc.BMC, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(
		context.TODO(), 3*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
			powerState, err := bmcClient.SystemPowerState()
			if err != nil {
				glog.V(tsparams.LogLevel).Infof("Failed to get system power state: %v", err)

				return false, err
			}

			if powerState != "Off" {
				glog.V(tsparams.LogLevel).Infof("System power state is not Off: %s", powerState)

				return false, nil
			}

			return true, nil
		})
}
