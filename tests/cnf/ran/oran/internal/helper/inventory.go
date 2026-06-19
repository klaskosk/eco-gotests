package helper

import (
	"context"
	"fmt"
	"slices"
	"time"

	bmhv1alpha1 "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/bmh"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/clients"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/schemes/ocm/clusterv1"
	inventoryv1alpha1 "github.com/openshift-kni/oran-o2ims/api/inventory/v1alpha1"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	goclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	resourcePoolNameLabel   = "resources.clcm.openshift.io/resourcePoolName"
	localManagedClusterName = "local-cluster"
)

var qualifyingProvisioningStates = []bmhv1alpha1.ProvisioningState{
	bmhv1alpha1.StateAvailable,
	bmhv1alpha1.StateProvisioning,
	bmhv1alpha1.StateProvisioned,
	bmhv1alpha1.StateExternallyProvisioned,
	bmhv1alpha1.StateDeprovisioning,
}

// ListLocationCRs lists Location CRs in the given namespace.
func ListLocationCRs(client *clients.Settings, namespace string) ([]inventoryv1alpha1.Location, error) {
	if err := client.AttachScheme(inventoryv1alpha1.AddToScheme); err != nil {
		return nil, fmt.Errorf("failed to attach Location scheme: %w", err)
	}

	var locationList inventoryv1alpha1.LocationList

	err := client.List(context.TODO(), &locationList, &goclient.ListOptions{Namespace: namespace})
	if err != nil {
		return nil, fmt.Errorf("failed to list Location CRs: %w", err)
	}

	return locationList.Items, nil
}

// ListOCloudSiteCRs lists OCloudSite CRs in the given namespace.
func ListOCloudSiteCRs(client *clients.Settings, namespace string) ([]inventoryv1alpha1.OCloudSite, error) {
	if err := client.AttachScheme(inventoryv1alpha1.AddToScheme); err != nil {
		return nil, fmt.Errorf("failed to attach OCloudSite scheme: %w", err)
	}

	var siteList inventoryv1alpha1.OCloudSiteList

	err := client.List(context.TODO(), &siteList, &goclient.ListOptions{Namespace: namespace})
	if err != nil {
		return nil, fmt.Errorf("failed to list OCloudSite CRs: %w", err)
	}

	return siteList.Items, nil
}

// ListResourcePoolCRs lists ResourcePool CRs in the given namespace.
func ListResourcePoolCRs(client *clients.Settings, namespace string) ([]inventoryv1alpha1.ResourcePool, error) {
	if err := client.AttachScheme(inventoryv1alpha1.AddToScheme); err != nil {
		return nil, fmt.Errorf("failed to attach ResourcePool scheme: %w", err)
	}

	var poolList inventoryv1alpha1.ResourcePoolList

	err := client.List(context.TODO(), &poolList, &goclient.ListOptions{Namespace: namespace})
	if err != nil {
		return nil, fmt.Errorf("failed to list ResourcePool CRs: %w", err)
	}

	return poolList.Items, nil
}

// ListManagedClusters lists ManagedCluster resources on the hub cluster, excluding the local hub cluster.
func ListManagedClusters(client *clients.Settings) ([]clusterv1.ManagedCluster, error) {
	if err := client.AttachScheme(clusterv1.Install); err != nil {
		return nil, fmt.Errorf("failed to attach ManagedCluster scheme: %w", err)
	}

	var managedClusterList clusterv1.ManagedClusterList

	err := client.List(context.TODO(), &managedClusterList)
	if err != nil {
		return nil, fmt.Errorf("failed to list ManagedClusters: %w", err)
	}

	var clusters []clusterv1.ManagedCluster

	for _, cluster := range managedClusterList.Items {
		if cluster.Name == localManagedClusterName {
			continue
		}

		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

// ListQualifiedBareMetalHostsForPool lists BareMetalHosts labeled for the given resource pool name that have
// completed hardware inspection and are in an available or later provisioning state.
func ListQualifiedBareMetalHostsForPool(
	client *clients.Settings, resourcePoolName string) ([]*bmh.BmhBuilder, error) {
	bmhList, err := bmh.ListInAllNamespaces(client, goclient.ListOptions{
		LabelSelector: labels.Set{resourcePoolNameLabel: resourcePoolName}.AsSelector(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list BareMetalHosts for resource pool %q: %w", resourcePoolName, err)
	}

	var qualified []*bmh.BmhBuilder

	for _, host := range bmhList {
		if host.Object.Status.OperationalStatus != bmhv1alpha1.OperationalStatusOK {
			continue
		}

		if !slices.Contains(qualifyingProvisioningStates, host.Object.Status.Provisioning.State) {
			continue
		}

		qualified = append(qualified, host)
	}

	return qualified, nil
}

// CreateTestLocation creates a Location CR for inventory subscription testing.
func CreateTestLocation(client *clients.Settings, name, description string) (*inventoryv1alpha1.Location, error) {
	if err := client.AttachScheme(inventoryv1alpha1.AddToScheme); err != nil {
		return nil, fmt.Errorf("failed to attach Location scheme: %w", err)
	}

	location := &inventoryv1alpha1.Location{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: tsparams.O2IMSNamespace,
		},
		Spec: inventoryv1alpha1.LocationSpec{
			Description: description,
			Address:     ptr.To("inventory test location"),
		},
	}

	err := client.Create(context.TODO(), location)
	if err != nil {
		return nil, fmt.Errorf("failed to create test Location CR %q: %w", name, err)
	}

	return location, nil
}

// DeleteTestLocation deletes a Location CR if it exists.
func DeleteTestLocation(client *clients.Settings, name string) error {
	if err := client.AttachScheme(inventoryv1alpha1.AddToScheme); err != nil {
		return fmt.Errorf("failed to attach Location scheme: %w", err)
	}

	location := &inventoryv1alpha1.Location{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: tsparams.O2IMSNamespace,
		},
	}

	err := client.Delete(context.TODO(), location)
	if err != nil {
		if goclient.IgnoreNotFound(err) != nil {
			return fmt.Errorf("failed to delete test Location CR %q: %w", name, err)
		}
	}

	return nil
}

// WaitForLocationInAPI waits until the given location is visible in the inventory API.
func WaitForLocationInAPI(
	inventoryClient *oranapi.InventoryClient, globalLocationID string, timeout time.Duration) (*oranapi.LocationInfo, error) {
	var location *oranapi.LocationInfo

	err := wait.PollUntilContextTimeout(
		context.TODO(), 3*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
			foundLocation, err := inventoryClient.GetLocation(globalLocationID)
			if err != nil {
				klog.V(tsparams.LogLevel).Infof("Location %q not yet visible in API: %v", globalLocationID, err)

				return false, nil
			}

			location = foundLocation

			return true, nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to wait for location %q in API: %w", globalLocationID, err)
	}

	return location, nil
}

// FindDeploymentManagerByName returns the deployment manager with the given name from a list response.
func FindDeploymentManagerByName(
	deploymentManagers []oranapi.DeploymentManager, name string) (*oranapi.DeploymentManager, bool) {
	for _, deploymentManager := range deploymentManagers {
		if deploymentManager.Name == name {
			return &deploymentManager, true
		}
	}

	return nil, false
}

// FindResourceByDescription returns the resource with the given description from a list response.
func FindResourceByDescription(resources []oranapi.Resource, description string) (*oranapi.Resource, bool) {
	for _, resource := range resources {
		if resource.Description == description {
			return &resource, true
		}
	}

	return nil, false
}

// BareMetalHostName returns the BareMetalHost metadata.name.
func BareMetalHostName(host *bmh.BmhBuilder) string {
	return host.Object.Name
}

// ResourcePoolLabelSelector returns a label selector for BareMetalHosts in the given resource pool.
func ResourcePoolLabelSelector(resourcePoolName string) labels.Set {
	return labels.Set{resourcePoolNameLabel: resourcePoolName}
}
