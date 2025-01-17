package helper

import (
	"fmt"

	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	pluginv1alpha1 "github.com/openshift-kni/oran-hwmgr-plugin/api/hwmgr-plugin/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// NewProvisioningRequest creates a ProvisioningRequest builder using the provided values. It sets the
// templateParameters automatically.
func NewProvisioningRequest(
	client *clients.Settings, clusterName, hostName, templateVersion string) *oran.ProvisioningRequestBuilder {
	prBuilder := oran.NewPRBuilder(client, clusterName, tsparams.ClusterTemplateName, templateVersion).
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
