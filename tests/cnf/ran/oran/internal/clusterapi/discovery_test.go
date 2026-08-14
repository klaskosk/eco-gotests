package clusterapi

import (
	"testing"

	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"github.com/stretchr/testify/assert"
)

func TestNodeClusterTypeKeyName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      NodeClusterTypeKey
		expected string
	}{
		{
			name:     "managed cluster",
			key:      NodeClusterTypeKey{Model: tsparams.ClusterModelManagedCluster, Vendor: "redhat", Version: "4.16"},
			expected: "managed-cluster-redhat-4.16",
		},
		{
			name:     "hub cluster",
			key:      NodeClusterTypeKey{Model: tsparams.ClusterModelHubCluster, Vendor: "redhat", Version: "4.17"},
			expected: "hub-cluster-redhat-4.17",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expected, testCase.key.Name())
		})
	}
}

func TestClusterResourceTypeKeyName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		key      ClusterResourceTypeKey
		expected string
	}{
		{
			name:     "uppercase architecture",
			key:      ClusterResourceTypeKey{Architecture: "x86_64", Cores: 96},
			expected: "X86_64 CPU with 96 Cores",
		},
		{
			name:     "mixed case architecture",
			key:      ClusterResourceTypeKey{Architecture: "Arm64", Cores: 32},
			expected: "ARM64 CPU with 32 Cores",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expected, testCase.key.Name())
		})
	}
}
