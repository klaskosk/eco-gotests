package clusterapi

import (
	"testing"

	"github.com/google/uuid"
	agentInstallV1Beta1 "github.com/rh-ecosystem-edge/eco-goinfra/pkg/schemes/assisted/api/v1beta1"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNodeClusterTypeKeyName(t *testing.T) {
	t.Parallel()

	key := NodeClusterTypeKey{Model: "managed-cluster", Vendor: "redhat", Version: "4.16"}
	assert.Equal(t, "managed-cluster-redhat-4.16", key.Name())
}

func TestClusterResourceTypeKeyName(t *testing.T) {
	t.Parallel()

	key := ClusterResourceTypeKey{Architecture: "x86_64", Cores: 96}
	assert.Equal(t, "X86_64 CPU with 96 Cores", key.Name())
}

func TestExpectedInventoryResourceID(t *testing.T) {
	t.Parallel()

	hwID := uuid.MustParse("a1b2c3d4-e5f6-7890-1234-567890abcdef")
	agent := &agentInstallV1Beta1.Agent{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "cluster",
			Name:      "host-0",
			Labels:    map[string]string{tsparams.HardwareManagerNodeIDLabel: hwID.String()},
		},
	}

	parsed, err := ExpectedInventoryResourceID(agent)
	require.NoError(t, err)
	assert.Equal(t, hwID, parsed)

	agentNoLabel := &agentInstallV1Beta1.Agent{ObjectMeta: metav1.ObjectMeta{Name: "host-1"}}
	parsed, err = ExpectedInventoryResourceID(agentNoLabel)
	require.NoError(t, err)
	assert.Equal(t, uuid.Nil, parsed)
}

func TestExpectedClusterResourceName(t *testing.T) {
	t.Parallel()

	agent := &agentInstallV1Beta1.Agent{
		Spec: agentInstallV1Beta1.AgentSpec{Hostname: "worker-0"},
	}
	assert.Equal(t, "worker-0", ExpectedClusterResourceName(agent))
}
