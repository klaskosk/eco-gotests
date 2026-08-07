package clusterapi

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/ocm"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	agentInstallV1Beta1 "github.com/rh-ecosystem-edge/eco-goinfra/pkg/schemes/assisted/api/v1beta1"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
)

const (
	expectedAPIVersion = "1.0.0"
	expectedURIPrefix  = "/o2ims-infrastructureCluster/v1"
)

// VerifyAPIVersions checks that the cluster API version response matches expected values.
func VerifyAPIVersions(versions oranapi.APIVersions) error {
	var errs []error

	apiVersions := derefSlice(versions.ApiVersions)
	if len(apiVersions) == 0 || apiVersions[0].Version == nil {
		errs = append(errs, fmt.Errorf("apiVersions[0].version: want non-nil, got nil"))
	} else {
		errs = appendMismatch(errs, "apiVersions[0].version", expectedAPIVersion, *apiVersions[0].Version)
	}

	if versions.UriPrefix == nil {
		errs = append(errs, fmt.Errorf("uriPrefix: want non-nil, got nil"))
	} else {
		errs = appendMismatch(errs, "uriPrefix", expectedURIPrefix, *versions.UriPrefix)
	}

	return errors.Join(errs...)
}

// VerifyNodeClusterTypeMatchesKey checks that an API NodeClusterType matches an expected type key.
func VerifyNodeClusterTypeMatchesKey(apiType oranapi.NodeClusterType, key NodeClusterTypeKey) error {
	var errs []error

	errs = appendMismatch(errs, "name", key.Name(), apiType.Name)
	errs = appendMismatch(errs, "description", key.Name(), apiType.Description)

	extensions := map[string]any{}
	if apiType.Extensions != nil {
		extensions = *apiType.Extensions
	}

	errs = appendMismatch(errs, "extensions.vendor", key.Vendor, extensionString(extensions, tsparams.ClusterVendorLabel))
	errs = appendMismatch(errs, "extensions.version", key.Version, extensionString(extensions, tsparams.ClusterVersionExtension))
	errs = appendMismatch(errs, "extensions.model", key.Model, extensionString(extensions, tsparams.ClusterModelExtension))

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
	errs = appendError(errs, idErr)
	errs = appendMismatch(errs, "nodeClusterId", expectedID, apiCluster.NodeClusterId)
	errs = appendMismatch(errs, "name", cluster.Definition.Name, apiCluster.Name)
	errs = appendMismatch(errs, "description", cluster.Definition.Name, apiCluster.Description)

	key, keyErr := NodeClusterTypeKeyFromManagedCluster(cluster)
	errs = appendError(errs, keyErr)

	typeIdx := slices.IndexFunc(nodeClusterTypes, func(nodeType oranapi.NodeClusterType) bool {
		return nodeType.NodeClusterTypeId == apiCluster.NodeClusterTypeId
	})
	if typeIdx == -1 {
		errs = append(errs, fmt.Errorf("nodeClusterTypeId %s not found in NodeClusterType list",
			apiCluster.NodeClusterTypeId))
	} else {
		errs = appendMismatch(errs, "nodeClusterType.name", key.Name(), nodeClusterTypes[typeIdx].Name)
	}

	extensions := map[string]any{}
	if apiCluster.Extensions != nil {
		extensions = *apiCluster.Extensions
	}

	if cluster.Definition.Labels != nil {
		errs = appendMismatch(errs, "extensions.vendor",
			cluster.Definition.Labels[tsparams.ClusterVendorLabel],
			extensionString(extensions, tsparams.ClusterVendorLabel))
		errs = appendMismatch(errs, "extensions.openshiftVersion",
			cluster.Definition.Labels[tsparams.OpenshiftVersionLabel],
			extensionString(extensions, tsparams.OpenshiftVersionLabel))
	}

	errs = appendMismatch(errs, "extensions.model", key.Model, extensionString(extensions, tsparams.ClusterModelExtension))

	return errors.Join(errs...)
}

// VerifyClusterResourceIDsExist reports an error when any clusterResourceId is missing from the listed ClusterResources.
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

	errs = appendMismatch(errs, "name", key.Name(), apiType.Name)
	errs = appendMismatch(errs, "description", key.Name(), apiType.Description)

	return errors.Join(errs...)
}

// VerifyClusterResourceMatchesAgent checks that an API ClusterResource matches the corresponding Agent.
func VerifyClusterResourceMatchesAgent(
	apiResource oranapi.ClusterResource,
	agent *agentInstallV1Beta1.Agent,
) error {
	var errs []error

	errs = appendMismatch(errs, "name", ExpectedClusterResourceName(agent), apiResource.Name)
	errs = appendMismatch(errs, "description", ExpectedClusterResourceName(agent), apiResource.Description)

	expectedResourceID, idErr := ExpectedInventoryResourceID(agent)
	errs = appendError(errs, idErr)
	errs = appendMismatch(errs, "resourceId", expectedResourceID, apiResource.ResourceId)

	if apiResource.ClusterResourceTypeId == uuid.Nil {
		errs = append(errs, fmt.Errorf("clusterResourceTypeId: want non-nil UUID, got %s", apiResource.ClusterResourceTypeId))
	}

	extensions := map[string]any{}
	if apiResource.Extensions != nil {
		extensions = *apiResource.Extensions
	}

	errs = appendError(errs, verifyClusterResourceCPUExtensions(extensions, agent))
	errs = appendError(errs, verifyClusterResourceMemoryExtensions(extensions, agent))
	errs = appendMismatch(errs, "extensions.role", string(agent.Status.Role), extensionString(extensions, "role"))

	return errors.Join(errs...)
}

func verifyClusterResourceCPUExtensions(extensions map[string]any, agent *agentInstallV1Beta1.Agent) error {
	cpuRaw, ok := extensions["cpu"]
	if !ok {
		return fmt.Errorf("extensions.cpu: want non-nil, got missing")
	}

	cpuMap, ok := asStringKeyedMap(cpuRaw)
	if !ok {
		return fmt.Errorf("extensions.cpu: want map, got %T", cpuRaw)
	}

	var errs []error

	errs = appendMismatch(errs, "extensions.cpu.architecture",
		agent.Status.Inventory.Cpu.Architecture, cpuMap["architecture"])
	errs = appendMismatch(errs, "extensions.cpu.cores",
		strconv.FormatInt(agent.Status.Inventory.Cpu.Count, 10), cpuMap["cores"])
	errs = appendMismatch(errs, "extensions.cpu.model",
		agent.Status.Inventory.Cpu.ModelName, cpuMap["model"])

	return errors.Join(errs...)
}

func verifyClusterResourceMemoryExtensions(extensions map[string]any, agent *agentInstallV1Beta1.Agent) error {
	memoryRaw, ok := extensions["memory"]
	if !ok {
		return fmt.Errorf("extensions.memory: want non-nil, got missing")
	}

	memoryMap, ok := asStringKeyedMap(memoryRaw)
	if !ok {
		return fmt.Errorf("extensions.memory: want map, got %T", memoryRaw)
	}

	expectedGiB := strconv.FormatInt(agent.Status.Inventory.Memory.PhysicalBytes/(1024*1024*1024), 10)

	return errors.Join(appendMismatch(nil, "extensions.memory.GiB", expectedGiB, memoryMap["GiB"])...)
}

// VerifyAlarmDictionaryStructure checks that an AlarmDictionary has the expected top-level fields and definitions.
func VerifyAlarmDictionaryStructure(dictionary oranapi.AlarmDictionary) error {
	var errs []error

	if dictionary.AlarmDictionaryId == uuid.Nil {
		errs = append(errs, fmt.Errorf("alarmDictionaryId: want non-nil UUID, got %s", dictionary.AlarmDictionaryId))
	}

	if dictionary.AlarmDictionaryVersion == "" {
		errs = append(errs, fmt.Errorf("alarmDictionaryVersion: want non-empty"))
	}

	if dictionary.EntityType == "" {
		errs = append(errs, fmt.Errorf("entityType: want non-empty"))
	}

	if dictionary.Vendor == "" {
		errs = append(errs, fmt.Errorf("vendor: want non-empty"))
	}

	if len(dictionary.AlarmDefinition) == 0 {
		errs = append(errs, fmt.Errorf("alarmDefinition: want non-empty"))
	}

	for i, definition := range dictionary.AlarmDefinition {
		if definition.AlarmName == "" {
			errs = append(errs, fmt.Errorf("alarmDefinition[%d].alarmName: want non-empty", i))
		}

		if definition.AlarmDescription == "" {
			errs = append(errs, fmt.Errorf("alarmDefinition[%d].alarmDescription: want non-empty", i))
		}

		// Severity lives in additionalFields; require the key but allow empty values when a Prometheus
		// rule omits the severity label.
		if definition.AlarmAdditionalFields == nil {
			errs = append(errs, fmt.Errorf(
				"alarmDefinition[%d].alarmAdditionalFields: want non-nil with %s key",
				i, tsparams.AlarmDefinitionSeverityField))

			continue
		}

		if _, ok := (*definition.AlarmAdditionalFields)[tsparams.AlarmDefinitionSeverityField]; !ok {
			errs = append(errs, fmt.Errorf(
				"alarmDefinition[%d].alarmAdditionalFields.%s: want key present",
				i, tsparams.AlarmDefinitionSeverityField))
		}
	}

	return errors.Join(errs...)
}

// NotificationRefersToNodeCluster reports whether the notification references the given NodeCluster.
func NotificationRefersToNodeCluster(
	notification *oranapi.ClusterChangeNotification,
	nodeClusterID uuid.UUID,
	nodeClusterName string,
) bool {
	if notification.ObjectRef != nil && strings.Contains(*notification.ObjectRef, nodeClusterID.String()) {
		return true
	}

	if notification.PostObjectState != nil {
		if value, ok := (*notification.PostObjectState)["nodeClusterId"]; ok && fmt.Sprint(value) == nodeClusterID.String() {
			return true
		}

		if value, ok := (*notification.PostObjectState)["name"]; ok && fmt.Sprint(value) == nodeClusterName {
			return true
		}
	}

	if notification.PriorObjectState != nil {
		if value, ok := (*notification.PriorObjectState)["nodeClusterId"]; ok && fmt.Sprint(value) == nodeClusterID.String() {
			return true
		}

		if value, ok := (*notification.PriorObjectState)["name"]; ok && fmt.Sprint(value) == nodeClusterName {
			return true
		}
	}

	return false
}

// NotificationHasExtensionLabel reports whether postObjectState.extensions contains the given label key/value.
func NotificationHasExtensionLabel(
	notification *oranapi.ClusterChangeNotification,
	labelKey, labelValue string,
) bool {
	if notification.PostObjectState == nil {
		return false
	}

	extensionsRaw, ok := (*notification.PostObjectState)["extensions"]
	if !ok || extensionsRaw == nil {
		return false
	}

	extensions, ok := asStringKeyedMap(extensionsRaw)
	if !ok {
		return false
	}

	return extensions[labelKey] == labelValue
}
