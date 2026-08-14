package clusterapi

import (
	"testing"

	"github.com/google/uuid"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen // Table-driven cases stay inline for readability.
func TestVerifyNodeClusterTypeMatchesKey(t *testing.T) {
	t.Parallel()

	key := NodeClusterTypeKey{Model: tsparams.ClusterModelHubCluster, Vendor: "redhat", Version: "4.16"}
	matchingExtensions := map[string]any{
		tsparams.ClusterVendorLabel:      "redhat",
		tsparams.ClusterVersionExtension: "4.16",
		tsparams.ClusterModelExtension:   tsparams.ClusterModelHubCluster,
	}

	tests := []struct {
		name    string
		apiType oranapi.NodeClusterType
		wantErr bool
	}{
		{
			name: "matching type",
			apiType: oranapi.NodeClusterType{
				Name:        key.Name(),
				Description: key.Name(),
				Extensions:  &matchingExtensions,
			},
		},
		{
			name: "nil extensions",
			apiType: oranapi.NodeClusterType{
				Name:        key.Name(),
				Description: key.Name(),
			},
			wantErr: true,
		},
		{
			name: "wrong name",
			apiType: oranapi.NodeClusterType{
				Name:        "wrong",
				Description: key.Name(),
				Extensions:  &matchingExtensions,
			},
			wantErr: true,
		},
		{
			name: "wrong vendor",
			apiType: oranapi.NodeClusterType{
				Name:        key.Name(),
				Description: key.Name(),
				Extensions: new(map[string]any{
					tsparams.ClusterVendorLabel:      "acme",
					tsparams.ClusterVersionExtension: "4.16",
					tsparams.ClusterModelExtension:   tsparams.ClusterModelHubCluster,
				}),
			},
			wantErr: true,
		},
		{
			name: "wrong version",
			apiType: oranapi.NodeClusterType{
				Name:        key.Name(),
				Description: key.Name(),
				Extensions: new(map[string]any{
					tsparams.ClusterVendorLabel:      "redhat",
					tsparams.ClusterVersionExtension: "4.17",
					tsparams.ClusterModelExtension:   tsparams.ClusterModelHubCluster,
				}),
			},
			wantErr: true,
		},
		{
			name: "wrong model",
			apiType: oranapi.NodeClusterType{
				Name:        key.Name(),
				Description: key.Name(),
				Extensions: new(map[string]any{
					tsparams.ClusterVendorLabel:      "redhat",
					tsparams.ClusterVersionExtension: "4.16",
					tsparams.ClusterModelExtension:   tsparams.ClusterModelManagedCluster,
				}),
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := VerifyNodeClusterTypeMatchesKey(testCase.apiType, key)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerifyClusterResourceTypeMatchesKey(t *testing.T) {
	t.Parallel()

	key := ClusterResourceTypeKey{Architecture: "x86_64", Cores: 96}

	tests := []struct {
		name    string
		apiType oranapi.ClusterResourceType
		wantErr bool
	}{
		{
			name:    "matching type",
			apiType: oranapi.ClusterResourceType{Name: key.Name(), Description: key.Name()},
		},
		{
			name:    "wrong name",
			apiType: oranapi.ClusterResourceType{Name: "wrong", Description: key.Name()},
			wantErr: true,
		},
		{
			name:    "wrong description",
			apiType: oranapi.ClusterResourceType{Name: key.Name(), Description: "wrong"},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := VerifyClusterResourceTypeMatchesKey(testCase.apiType, key)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerifyClusterResourceIDsExist(t *testing.T) {
	t.Parallel()

	resourceID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	resources := []oranapi.ClusterResource{{ClusterResourceId: resourceID}}

	tests := []struct {
		name      string
		ids       []uuid.UUID
		resources []oranapi.ClusterResource
		wantErr   bool
	}{
		{
			name:      "all present",
			ids:       []uuid.UUID{resourceID},
			resources: resources,
		},
		{
			name:      "empty ids",
			resources: resources,
		},
		{
			name: "missing id",
			ids: []uuid.UUID{
				resourceID,
				uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			},
			resources: resources,
			wantErr:   true,
		},
		{
			name:    "empty resources",
			ids:     []uuid.UUID{resourceID},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := VerifyClusterResourceIDsExist(testCase.ids, testCase.resources)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
