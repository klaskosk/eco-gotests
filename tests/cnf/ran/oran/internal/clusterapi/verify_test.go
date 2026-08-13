package clusterapi

import (
	"testing"

	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"github.com/stretchr/testify/assert"
)

func TestVerifyNodeClusterTypeMatchesKey(t *testing.T) {
	t.Parallel()

	key := NodeClusterTypeKey{Model: tsparams.ClusterModelHubCluster, Vendor: "redhat", Version: "4.16"}
	extensions := map[string]any{
		tsparams.ClusterVendorLabel:      "redhat",
		tsparams.ClusterVersionExtension: "4.16",
		tsparams.ClusterModelExtension:   tsparams.ClusterModelHubCluster,
	}

	apiType := oranapi.NodeClusterType{
		Name:        key.Name(),
		Description: key.Name(),
		Extensions:  &extensions,
	}

	assert.NoError(t, VerifyNodeClusterTypeMatchesKey(apiType, key))

	apiType.Extensions = nil
	err := VerifyNodeClusterTypeMatchesKey(apiType, key)
	assert.Error(t, err)
}

func TestVerifyClusterResourceTypeMatchesKey(t *testing.T) {
	t.Parallel()

	key := ClusterResourceTypeKey{Architecture: "x86_64", Cores: 96}
	apiType := oranapi.ClusterResourceType{Name: key.Name(), Description: key.Name()}

	assert.NoError(t, VerifyClusterResourceTypeMatchesKey(apiType, key))

	err := VerifyClusterResourceTypeMatchesKey(
		oranapi.ClusterResourceType{Name: "wrong", Description: "wrong"}, key)
	assert.Error(t, err)
}
