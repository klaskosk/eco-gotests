package o2imstest

import (
	"testing"

	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"github.com/stretchr/testify/assert"
)

func TestAsStringKeyedMap(t *testing.T) {
	t.Parallel()

	fromAny, converted := AsStringKeyedMap(map[string]any{"cores": "96", "architecture": "x86_64"})
	assert.True(t, converted)
	assert.Equal(t, "96", fromAny["cores"])
	assert.Equal(t, "x86_64", fromAny["architecture"])

	fromString, converted := AsStringKeyedMap(map[string]string{"GiB": "256"})
	assert.True(t, converted)
	assert.Equal(t, "256", fromString["GiB"])

	_, converted = AsStringKeyedMap("not-a-map")
	assert.False(t, converted)
}

func TestVerifyAPIVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		versions          oranapi.APIVersions
		expectedVersion   string
		expectedURIPrefix string
		wantErr           bool
	}{
		{
			name: "valid cluster versions",
			versions: oranapi.APIVersions{
				ApiVersions: &[]oranapi.APIVersion{{Version: new(tsparams.ClusterAPIVersion)}},
				UriPrefix:   new(tsparams.ClusterAPIURIPrefix),
			},
			expectedVersion:   tsparams.ClusterAPIVersion,
			expectedURIPrefix: tsparams.ClusterAPIURIPrefix,
		},
		{
			name: "valid inventory versions",
			versions: oranapi.APIVersions{
				ApiVersions: &[]oranapi.APIVersion{{Version: new(tsparams.InventoryAPIVersion)}},
				UriPrefix:   new(tsparams.InventoryAPIURIPrefix),
			},
			expectedVersion:   tsparams.InventoryAPIVersion,
			expectedURIPrefix: tsparams.InventoryAPIURIPrefix,
		},
		{
			name: "wrong version",
			versions: oranapi.APIVersions{
				ApiVersions: &[]oranapi.APIVersion{{Version: new(tsparams.InventoryAPIVersion)}},
				UriPrefix:   new(tsparams.ClusterAPIURIPrefix),
			},
			expectedVersion:   tsparams.ClusterAPIVersion,
			expectedURIPrefix: tsparams.ClusterAPIURIPrefix,
			wantErr:           true,
		},
		{
			name: "wrong uriPrefix",
			versions: oranapi.APIVersions{
				ApiVersions: &[]oranapi.APIVersion{{Version: new(tsparams.ClusterAPIVersion)}},
				UriPrefix:   new("/wrong"),
			},
			expectedVersion:   tsparams.ClusterAPIVersion,
			expectedURIPrefix: tsparams.ClusterAPIURIPrefix,
			wantErr:           true,
		},
		{
			name:              "missing fields",
			versions:          oranapi.APIVersions{},
			expectedVersion:   tsparams.ClusterAPIVersion,
			expectedURIPrefix: tsparams.ClusterAPIURIPrefix,
			wantErr:           true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := VerifyAPIVersions(testCase.versions, testCase.expectedVersion, testCase.expectedURIPrefix)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRefersTo(t *testing.T) {
	t.Parallel()

	const (
		objectID = "pool-1"
		name     = "test-pool"
	)

	post := map[string]any{"resourcePoolId": objectID, "name": name}
	prior := map[string]any{"resourcePoolId": objectID}

	assert.True(t, RefersTo(new("/resourcePools/"+objectID), nil, nil, "resourcePoolId", objectID, "name", name))
	assert.True(t, RefersTo(nil, &post, nil, "resourcePoolId", objectID, "name", name))
	assert.True(t, RefersTo(nil, nil, &prior, "resourcePoolId", objectID, "name", name))
	assert.True(t, RefersTo(nil, &map[string]any{"name": name}, nil, "resourcePoolId", objectID, "name", name))
	assert.False(t, RefersTo(new("/resourcePools/other"), nil, nil, "resourcePoolId", objectID, "name", name))
	assert.False(t, RefersTo(nil, nil, nil, "resourcePoolId", objectID, "name", name))
}

func TestExtensionString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "", ExtensionString(nil, tsparams.ClusterModelExtension))

	extensions := map[string]any{tsparams.ClusterModelExtension: tsparams.ClusterModelHubCluster}
	assert.Equal(t, tsparams.ClusterModelHubCluster, ExtensionString(extensions, tsparams.ClusterModelExtension))
	assert.Equal(t, "", ExtensionString(extensions, "missing"))
}
