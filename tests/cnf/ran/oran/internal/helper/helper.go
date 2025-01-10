package helper

import (
	"encoding/json"

	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"k8s.io/apimachinery/pkg/runtime"
)

// NewProvisioningRequest creates a ProvisioningRequest builder using the provided values. It sets the
// templateParameters automatically.
func NewProvisioningRequest(
	client *clients.Settings, clusterName, hostName, templateVersion string) *oran.ProvisioningRequestBuilder {
	prBuilder := oran.NewPRBuilder(client, clusterName, tsparams.ClusterTemplateName, templateVersion)

	templateParams := map[string]any{
		"nodeClusterName":          clusterName,
		"oCloudSiteId":             clusterName,
		"policyTemplateParameters": map[string]any{},
		"clusterInstanceParameters": map[string]any{
			"clusterName": clusterName,
			"nodes": []map[string]any{{
				"hostName": hostName,
			}},
		},
	}
	templateMarshaled, err := json.Marshal(templateParams)
	templateRaw := runtime.RawExtension{Raw: templateMarshaled}

	if err != nil {
		// TODO: add actual error handling
		panic(err)
	}

	prBuilder.Definition.Spec.TemplateParameters = templateRaw

	return prBuilder
}
