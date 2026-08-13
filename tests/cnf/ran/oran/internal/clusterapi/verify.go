package clusterapi

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/google/uuid"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/ocm"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	agentInstallV1Beta1 "github.com/rh-ecosystem-edge/eco-goinfra/pkg/schemes/assisted/api/v1beta1"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/o2imstest"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
)

// VerifyNodeClusterTypeMatchesKey checks that an API NodeClusterType matches an expected type key.
func VerifyNodeClusterTypeMatchesKey(apiType oranapi.NodeClusterType, key NodeClusterTypeKey) error {
	var errs []error

	errs = o2imstest.AppendMismatch(errs, "name", key.Name(), apiType.Name)
	errs = o2imstest.AppendMismatch(errs, "description", key.Name(), apiType.Description)

	extensions := map[string]any{}
	if apiType.Extensions != nil {
		extensions = *apiType.Extensions
	}

	errs = o2imstest.AppendMismatch(errs, "extensions.vendor", key.Vendor,
		o2imstest.ExtensionString(extensions, tsparams.ClusterVendorLabel))
	errs = o2imstest.AppendMismatch(errs, "extensions.version", key.Version,
		o2imstest.ExtensionString(extensions, tsparams.ClusterVersionExtension))
	errs = o2imstest.AppendMismatch(errs, "extensions.model", key.Model,
		o2imstest.ExtensionString(extensions, tsparams.ClusterModelExtension))

	return errors.Join(errs...)
}

// VerifyNodeClusterMatchesManagedCluster checks that an API NodeCluster matches a ManagedCluster.
func VerifyNodeClusterMatchesManagedCluster(
	apiCluster oranapi.NodeCluster,
	cluster *ocm.ManagedClusterBuilder,
	nodeClusterTypes []oranapi.NodeClusterType,
) error {
	var errs []error

	expectedID, idErr := ManagedClusterID(cluster)
	errs = o2imstest.AppendError(errs, idErr)
	errs = o2imstest.AppendMismatch(errs, "nodeClusterId", expectedID, apiCluster.NodeClusterId)
	errs = o2imstest.AppendMismatch(errs, "name", cluster.Definition.Name, apiCluster.Name)
	errs = o2imstest.AppendMismatch(errs, "description", cluster.Definition.Name, apiCluster.Description)

	key, keyErr := NodeClusterTypeKeyFromManagedCluster(cluster)
	errs = o2imstest.AppendError(errs, keyErr)

	typeIdx := slices.IndexFunc(nodeClusterTypes, func(nodeType oranapi.NodeClusterType) bool {
		return nodeType.NodeClusterTypeId == apiCluster.NodeClusterTypeId
	})
	if typeIdx == -1 {
		errs = append(errs, fmt.Errorf("nodeClusterTypeId %s not found in NodeClusterType list",
			apiCluster.NodeClusterTypeId))
	} else {
		errs = o2imstest.AppendMismatch(errs, "nodeClusterType.name", key.Name(), nodeClusterTypes[typeIdx].Name)
	}

	extensions := map[string]any{}
	if apiCluster.Extensions != nil {
		extensions = *apiCluster.Extensions
	}

	if cluster.Definition.Labels != nil {
		errs = o2imstest.AppendMismatch(errs, "extensions.vendor",
			cluster.Definition.Labels[tsparams.ClusterVendorLabel],
			o2imstest.ExtensionString(extensions, tsparams.ClusterVendorLabel))
		errs = o2imstest.AppendMismatch(errs, "extensions.openshiftVersion",
			cluster.Definition.Labels[tsparams.OpenshiftVersionLabel],
			o2imstest.ExtensionString(extensions, tsparams.OpenshiftVersionLabel))
	}

	errs = o2imstest.AppendMismatch(errs, "extensions.model", key.Model,
		o2imstest.ExtensionString(extensions, tsparams.ClusterModelExtension))

	return errors.Join(errs...)
}

// VerifyClusterResourceIDsExist reports an error when any clusterResourceId is missing from the
// listed ClusterResources.
func VerifyClusterResourceIDsExist(ids []uuid.UUID, resources []oranapi.ClusterResource) error {
	resourceIDs := make([]string, 0, len(resources))
	for _, resource := range resources {
		resourceIDs = append(resourceIDs, resource.ClusterResourceId.String())
	}

	var errs []error

	for _, id := range ids {
		if !slices.Contains(resourceIDs, id.String()) {
			errs = append(errs, fmt.Errorf("clusterResourceId %s not found in ClusterResource list", id))
		}
	}

	return errors.Join(errs...)
}

// VerifyClusterResourceTypeMatchesKey checks that an API ClusterResourceType matches an expected type key.
func VerifyClusterResourceTypeMatchesKey(apiType oranapi.ClusterResourceType, key ClusterResourceTypeKey) error {
	var errs []error

	errs = o2imstest.AppendMismatch(errs, "name", key.Name(), apiType.Name)
	errs = o2imstest.AppendMismatch(errs, "description", key.Name(), apiType.Description)

	return errors.Join(errs...)
}

// VerifyClusterResourceMatchesAgent checks that an API ClusterResource matches the corresponding Agent.
func VerifyClusterResourceMatchesAgent(
	apiResource oranapi.ClusterResource,
	agent *agentInstallV1Beta1.Agent,
) error {
	var errs []error

	errs = o2imstest.AppendMismatch(errs, "name", ExpectedClusterResourceName(agent), apiResource.Name)
	errs = o2imstest.AppendMismatch(errs, "description", ExpectedClusterResourceName(agent), apiResource.Description)

	expectedResourceID, idErr := ExpectedInventoryResourceID(agent)
	errs = o2imstest.AppendError(errs, idErr)
	errs = o2imstest.AppendMismatch(errs, "resourceId", expectedResourceID, apiResource.ResourceId)

	if apiResource.ClusterResourceTypeId == uuid.Nil {
		errs = append(errs, fmt.Errorf("clusterResourceTypeId: want non-nil UUID, got %s", apiResource.ClusterResourceTypeId))
	}

	extensions := map[string]any{}
	if apiResource.Extensions != nil {
		extensions = *apiResource.Extensions
	}

	errs = o2imstest.AppendError(errs, verifyClusterResourceCPUExtensions(extensions, agent))
	errs = o2imstest.AppendError(errs, verifyClusterResourceMemoryExtensions(extensions, agent))
	errs = o2imstest.AppendMismatch(errs, "extensions.role", string(agent.Status.Role),
		o2imstest.ExtensionString(extensions, "role"))

	return errors.Join(errs...)
}

// verifyClusterResourceCPUExtensions checks extensions.cpu against the Agent inventory CPU fields.
func verifyClusterResourceCPUExtensions(extensions map[string]any, agent *agentInstallV1Beta1.Agent) error {
	cpuRaw, found := extensions["cpu"]
	if !found {
		return fmt.Errorf("extensions.cpu: want non-nil, got missing")
	}

	cpuMap, converted := o2imstest.AsStringKeyedMap(cpuRaw)
	if !converted {
		return fmt.Errorf("extensions.cpu: want map, got %T", cpuRaw)
	}

	var errs []error

	errs = o2imstest.AppendMismatch(errs, "extensions.cpu.architecture",
		agent.Status.Inventory.Cpu.Architecture, cpuMap["architecture"])
	errs = o2imstest.AppendMismatch(errs, "extensions.cpu.cores",
		strconv.FormatInt(agent.Status.Inventory.Cpu.Count, 10), cpuMap["cores"])
	errs = o2imstest.AppendMismatch(errs, "extensions.cpu.model",
		agent.Status.Inventory.Cpu.ModelName, cpuMap["model"])

	return errors.Join(errs...)
}

// verifyClusterResourceMemoryExtensions checks extensions.memory.GiB against the Agent inventory memory size.
func verifyClusterResourceMemoryExtensions(extensions map[string]any, agent *agentInstallV1Beta1.Agent) error {
	memoryRaw, found := extensions["memory"]
	if !found {
		return fmt.Errorf("extensions.memory: want non-nil, got missing")
	}

	memoryMap, converted := o2imstest.AsStringKeyedMap(memoryRaw)
	if !converted {
		return fmt.Errorf("extensions.memory: want map, got %T", memoryRaw)
	}

	expectedGiB := strconv.FormatInt(agent.Status.Inventory.Memory.PhysicalBytes/(1024*1024*1024), 10)
	if expectedGiB != memoryMap["GiB"] {
		return fmt.Errorf("extensions.memory.GiB: want %#v, got %#v", expectedGiB, memoryMap["GiB"])
	}

	return nil
}
