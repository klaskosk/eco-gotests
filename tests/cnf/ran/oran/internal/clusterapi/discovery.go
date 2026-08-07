package clusterapi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/assisted"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/clients"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/ocm"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	agentInstallV1Beta1 "github.com/rh-ecosystem-edge/eco-goinfra/pkg/schemes/assisted/api/v1beta1"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
)

// NodeClusterTypeKey identifies a distinct NodeClusterType derived from eligible ManagedClusters.
type NodeClusterTypeKey struct {
	Model   string
	Vendor  string
	Version string
}

// Name returns the NodeClusterType name pattern {model}-{vendor}-{version}.
func (key NodeClusterTypeKey) Name() string {
	return fmt.Sprintf("%s-%s-%s", key.Model, key.Vendor, key.Version)
}

// ClusterResourceTypeKey identifies a distinct ClusterResourceType derived from eligible Agents.
type ClusterResourceTypeKey struct {
	Architecture string
	Cores        int64
}

// Name returns the ClusterResourceType name pattern "{ARCH} CPU with {N} Cores".
func (key ClusterResourceTypeKey) Name() string {
	return fmt.Sprintf("%s CPU with %d Cores", strings.ToUpper(key.Architecture), key.Cores)
}

// IsHubManagedCluster reports whether cluster is the ACM hub (local-cluster=true).
func IsHubManagedCluster(cluster *ocm.ManagedClusterBuilder) bool {
	if cluster == nil || cluster.Definition == nil || cluster.Definition.Labels == nil {
		return false
	}

	value, found := cluster.Definition.Labels[tsparams.LocalClusterLabel]
	if !found {
		return false
	}

	localCluster, err := strconv.ParseBool(value)

	return err == nil && localCluster
}

// ManagedClusterModel returns hub-cluster or managed-cluster for the given ManagedCluster.
func ManagedClusterModel(cluster *ocm.ManagedClusterBuilder) string {
	if IsHubManagedCluster(cluster) {
		return tsparams.ClusterModelHubCluster
	}

	return tsparams.ClusterModelManagedCluster
}

// NodeClusterTypeKeyFromManagedCluster builds the type key from an eligible ManagedCluster's labels.
func NodeClusterTypeKeyFromManagedCluster(cluster *ocm.ManagedClusterBuilder) (NodeClusterTypeKey, error) {
	if cluster == nil || cluster.Definition == nil || cluster.Definition.Labels == nil {
		return NodeClusterTypeKey{}, fmt.Errorf("managed cluster or labels are nil")
	}

	vendor, found := cluster.Definition.Labels[tsparams.ClusterVendorLabel]
	if !found {
		return NodeClusterTypeKey{}, fmt.Errorf("managed cluster %s missing %s label",
			cluster.Definition.Name, tsparams.ClusterVendorLabel)
	}

	version := cluster.Definition.Labels[tsparams.OpenshiftVersionLabel]

	return NodeClusterTypeKey{
		Model:   ManagedClusterModel(cluster),
		Vendor:  vendor,
		Version: version,
	}, nil
}

// CollectNodeClusterTypeKeys returns distinct NodeClusterType keys from eligible ManagedClusters.
func CollectNodeClusterTypeKeys(hubClient *clients.Settings) (map[NodeClusterTypeKey]struct{}, error) {
	eligible, err := ocm.ListNodeClusterEligibleManagedClusters(hubClient)
	if err != nil {
		return nil, fmt.Errorf("failed to list node-cluster-eligible ManagedClusters: %w", err)
	}

	keys := make(map[NodeClusterTypeKey]struct{}, len(eligible))

	for _, cluster := range eligible {
		key, keyErr := NodeClusterTypeKeyFromManagedCluster(cluster)
		if keyErr != nil {
			return nil, keyErr
		}

		keys[key] = struct{}{}
	}

	return keys, nil
}

// FindSpokeEligibleManagedCluster returns an eligible non-hub ManagedCluster, preferring preferredName when set.
func FindSpokeEligibleManagedCluster(
	hubClient *clients.Settings,
	preferredName string,
) (*ocm.ManagedClusterBuilder, error) {
	eligible, err := ocm.ListNodeClusterEligibleManagedClusters(hubClient)
	if err != nil {
		return nil, fmt.Errorf("failed to list node-cluster-eligible ManagedClusters: %w", err)
	}

	var fallback *ocm.ManagedClusterBuilder

	for _, cluster := range eligible {
		if IsHubManagedCluster(cluster) {
			continue
		}

		if preferredName != "" && cluster.Definition.Name == preferredName {
			return cluster, nil
		}

		if fallback == nil {
			fallback = cluster
		}
	}

	if fallback == nil {
		return nil, fmt.Errorf("no eligible spoke ManagedCluster found")
	}

	return fallback, nil
}

// ManagedClusterID returns the clusterID label UUID for a ManagedCluster.
func ManagedClusterID(cluster *ocm.ManagedClusterBuilder) (uuid.UUID, error) {
	if cluster == nil || cluster.Definition == nil || cluster.Definition.Labels == nil {
		return uuid.Nil, fmt.Errorf("managed cluster or labels are nil")
	}

	clusterID, found := cluster.Definition.Labels[tsparams.ClusterIDLabel]
	if !found {
		return uuid.Nil, fmt.Errorf("managed cluster %s missing %s label",
			cluster.Definition.Name, tsparams.ClusterIDLabel)
	}

	parsed, err := uuid.Parse(clusterID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("managed cluster %s has invalid %s label %q: %w",
			cluster.Definition.Name, tsparams.ClusterIDLabel, clusterID, err)
	}

	return parsed, nil
}

// CollectClusterResourceTypeKeys returns distinct ClusterResourceType keys from ORAN-eligible Agents.
func CollectClusterResourceTypeKeys(hubClient *clients.Settings) (map[ClusterResourceTypeKey]struct{}, error) {
	agents, err := assisted.ListORANEligibleAgents(hubClient)
	if err != nil {
		return nil, fmt.Errorf("failed to list ORAN-eligible Agents: %w", err)
	}

	keys := make(map[ClusterResourceTypeKey]struct{}, len(agents))

	for _, agent := range agents {
		key := ClusterResourceTypeKeyFromAgent(agent.Definition)
		keys[key] = struct{}{}
	}

	return keys, nil
}

// ClusterResourceTypeKeyFromAgent builds the type key from an Agent's inventory CPU fields.
func ClusterResourceTypeKeyFromAgent(agent *agentInstallV1Beta1.Agent) ClusterResourceTypeKey {
	return ClusterResourceTypeKey{
		Architecture: agent.Status.Inventory.Cpu.Architecture,
		Cores:        agent.Status.Inventory.Cpu.Count,
	}
}

// ExpectedClusterResourceName returns the ClusterResource name the API derives from an Agent.
func ExpectedClusterResourceName(agent *agentInstallV1Beta1.Agent) string {
	return agent.Spec.Hostname
}

// ExpectedAgentExternalID returns the namespace/name identity the collector uses as ExternalID.
// The public Cluster API does not expose ExternalID; this is used for stable error messages and matching.
func ExpectedAgentExternalID(agent *agentInstallV1Beta1.Agent) string {
	return fmt.Sprintf("%s/%s", agent.Namespace, agent.Name)
}

// ExpectedInventoryResourceID returns the ClusterResource.resourceId derived from the Agent hwMgrNodeId label.
func ExpectedInventoryResourceID(agent *agentInstallV1Beta1.Agent) (uuid.UUID, error) {
	if agent.Labels == nil {
		return uuid.Nil, nil
	}

	hwMgrNodeID, found := agent.Labels[tsparams.HardwareManagerNodeIDLabel]
	if !found {
		return uuid.Nil, nil
	}

	parsed, err := uuid.Parse(hwMgrNodeID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("agent %s/%s has invalid %s label %q: %w",
			agent.Namespace, agent.Name, tsparams.HardwareManagerNodeIDLabel, hwMgrNodeID, err)
	}

	return parsed, nil
}

// FindClusterResourceForAgent returns the ClusterResource corresponding to agent.
// Matching prefers resourceId from the hwMgrNodeId label when present; otherwise matches on
// hostname plus CPU architecture/cores so duplicate hostnames are less likely to collide.
func FindClusterResourceForAgent(
	resources []oranapi.ClusterResource,
	agent *agentInstallV1Beta1.Agent,
) (oranapi.ClusterResource, error) {
	expectedName := ExpectedClusterResourceName(agent)
	expectedResourceID, idErr := ExpectedInventoryResourceID(agent)
	if idErr != nil {
		return oranapi.ClusterResource{}, idErr
	}

	key := ClusterResourceTypeKeyFromAgent(agent)

	if expectedResourceID != uuid.Nil {
		for _, resource := range resources {
			if resource.ResourceId == expectedResourceID {
				return resource, nil
			}
		}
	}

	var match *oranapi.ClusterResource

	for i := range resources {
		resource := &resources[i]
		if resource.Name != expectedName {
			continue
		}

		if !clusterResourceCPUMatchesKey(*resource, key) {
			continue
		}

		if match != nil {
			return oranapi.ClusterResource{}, fmt.Errorf(
				"ambiguous ClusterResource match for Agent %s (name %s): multiple candidates",
				ExpectedAgentExternalID(agent), expectedName)
		}

		match = resource
	}

	if match == nil {
		return oranapi.ClusterResource{}, fmt.Errorf(
			"no ClusterResource matched Agent %s (name %s, resourceId %s)",
			ExpectedAgentExternalID(agent), expectedName, expectedResourceID)
	}

	return *match, nil
}

func clusterResourceCPUMatchesKey(resource oranapi.ClusterResource, key ClusterResourceTypeKey) bool {
	if resource.Extensions == nil {
		return false
	}

	cpuRaw, ok := (*resource.Extensions)["cpu"]
	if !ok {
		return false
	}

	cpuMap, ok := asStringKeyedMap(cpuRaw)
	if !ok {
		return false
	}

	return cpuMap["architecture"] == key.Architecture &&
		cpuMap["cores"] == strconv.FormatInt(key.Cores, 10)
}
