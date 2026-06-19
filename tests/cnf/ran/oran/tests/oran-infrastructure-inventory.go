package tests

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/bmh"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api/filter"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/reportxml"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/auth"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/helper"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	subscriber "github.com/rh-ecosystem-edge/eco-gotests/tests/internal/oran-subscriber"
	"k8s.io/utils/ptr"
)

var inventorySubscriberURL = "https://" + RANConfig.GetAppsURL(tsparams.SubscriberSubdomain)

const inventoryBMHPrerequisite = "requires a BareMetalHost labeled with " +
	"resources.clcm.openshift.io/resourcePoolName in available or later provisioning state"

var _ = Describe("ORAN Infrastructure Inventory Tests",
	Label(tsparams.LabelPreProvision, tsparams.LabelInfrastructureInventory), func() {
		var inventoryClient *oranapi.InventoryClient

		BeforeEach(func() {
			var err error

			By("creating the O2IMS inventory API client")

			clientBuilder, err := auth.NewClientBuilderForConfig(RANConfig)
			Expect(err).ToNot(HaveOccurred(), "Failed to create the O2IMS API client builder")

			inventoryClient, err = clientBuilder.BuildInventory()
			Expect(err).ToNot(HaveOccurred(), "Failed to create the O2IMS inventory API client")
		})

		// 00001 - Successfully retrieve O-Cloud metadata and API versions
		It("retrieves O-Cloud metadata and API versions", reportxml.ID("00001"), func() {
			By("getting all API versions")

			apiVersions, err := inventoryClient.GetAPIVersions()
			Expect(err).ToNot(HaveOccurred(), "Failed to get API versions")
			Expect(apiVersions.ApiVersions).ToNot(BeNil())
			Expect(*apiVersions.ApiVersions).ToNot(BeEmpty())
			Expect(apiVersions.UriPrefix).ToNot(BeNil())
			Expect(*apiVersions.UriPrefix).To(ContainSubstring("infrastructureInventory"))

			for _, version := range *apiVersions.ApiVersions {
				Expect(version.Version).ToNot(BeNil())
				Expect(*version.Version).ToNot(BeEmpty())
			}

			By("getting O-Cloud info")

			cloudInfo, err := inventoryClient.GetCloudInfo()
			Expect(err).ToNot(HaveOccurred(), "Failed to get O-Cloud info")
			Expect(cloudInfo.OCloudId).ToNot(Equal(uuid.Nil))
			Expect(cloudInfo.GlobalCloudId).ToNot(Equal(uuid.Nil))
			Expect(cloudInfo.Name).ToNot(BeEmpty())
			Expect(cloudInfo.Description).ToNot(BeEmpty())
			Expect(cloudInfo.ServiceUri).ToNot(BeEmpty())
			Expect(cloudInfo.Extensions).ToNot(BeNil())

			By("getting minor API versions for major version 1")

			minorVersions, err := inventoryClient.GetMinorAPIVersions()
			Expect(err).ToNot(HaveOccurred(), "Failed to get minor API versions")
			Expect(minorVersions.ApiVersions).ToNot(BeNil())
			Expect(*minorVersions.ApiVersions).ToNot(BeEmpty())
		})

		// 00002 - Successfully list and retrieve locations
		It("lists and retrieves locations", reportxml.ID("00002"), func() {
			By("listing Location CRs on the hub cluster")

			locationCRs, err := helper.ListLocationCRs(HubAPIClient, tsparams.O2IMSNamespace)
			Expect(err).ToNot(HaveOccurred(), "Failed to list Location CRs")
			Expect(locationCRs).ToNot(BeEmpty(), "At least one Location CR should be present")

			By("listing locations from the API")

			locations, err := inventoryClient.ListLocations()
			Expect(err).ToNot(HaveOccurred(), "Failed to list locations from the API")
			Expect(locations).To(HaveLen(len(locationCRs)), "API location count should match Location CR count")

			By("retrieving each location by globalLocationId")

			for _, locationCR := range locationCRs {
				location, err := inventoryClient.GetLocation(locationCR.Name)
				Expect(err).ToNot(HaveOccurred(), "Failed to get location %s", locationCR.Name)
				Expect(location.GlobalLocationId).To(Equal(locationCR.Name))
				Expect(location.Description).To(Equal(locationCR.Spec.Description))
			}
		})

		// 00003 - Successfully list and retrieve O-Cloud sites
		It("lists and retrieves O-Cloud sites", reportxml.ID("00003"), func() {
			By("listing OCloudSite CRs on the hub cluster")

			siteCRs, err := helper.ListOCloudSiteCRs(HubAPIClient, tsparams.O2IMSNamespace)
			Expect(err).ToNot(HaveOccurred(), "Failed to list OCloudSite CRs")
			Expect(siteCRs).ToNot(BeEmpty(), "At least one OCloudSite CR should be present")

			By("listing O-Cloud sites from the API")

			sites, err := inventoryClient.ListOCloudSites()
			Expect(err).ToNot(HaveOccurred(), "Failed to list O-Cloud sites from the API")
			Expect(sites).To(HaveLen(len(siteCRs)), "API site count should match OCloudSite CR count")

			By("retrieving each O-Cloud site by oCloudSiteId")

			for _, siteCR := range siteCRs {
				siteID, err := uuid.Parse(string(siteCR.UID))
				Expect(err).ToNot(HaveOccurred(), "Failed to parse OCloudSite UID")

				site, err := inventoryClient.GetOCloudSite(siteID)
				Expect(err).ToNot(HaveOccurred(), "Failed to get O-Cloud site %s", siteCR.Name)
				Expect(site.OCloudSiteId).To(Equal(siteID))
				Expect(site.Name).To(Equal(siteCR.Name))
			}
		})

		// 00004 - Successfully list and retrieve resource pools
		It("lists and retrieves resource pools", reportxml.ID("00004"), func() {
			By("listing ResourcePool CRs on the hub cluster")

			poolCRs, err := helper.ListResourcePoolCRs(HubAPIClient, tsparams.O2IMSNamespace)
			Expect(err).ToNot(HaveOccurred(), "Failed to list ResourcePool CRs")
			Expect(poolCRs).ToNot(BeEmpty(), "At least one ResourcePool CR should be present")

			By("listing resource pools from the API")

			pools, err := inventoryClient.ListResourcePools()
			Expect(err).ToNot(HaveOccurred(), "Failed to list resource pools from the API")
			Expect(pools).To(HaveLen(len(poolCRs)), "API pool count should match ResourcePool CR count")

			By("retrieving each resource pool by resourcePoolId")

			for _, poolCR := range poolCRs {
				poolID, err := uuid.Parse(string(poolCR.UID))
				Expect(err).ToNot(HaveOccurred(), "Failed to parse ResourcePool UID")

				pool, err := inventoryClient.GetResourcePool(poolID)
				Expect(err).ToNot(HaveOccurred(), "Failed to get resource pool %s", poolCR.Name)
				Expect(pool.ResourcePoolId).To(Equal(poolID))
				Expect(pool.Name).To(Equal(poolCR.Name))
			}
		})

		// 00005 - Successfully list and retrieve resources in a resource pool
		It("lists and retrieves resources in a resource pool", reportxml.ID("00005"), func() {
			By("identifying a ResourcePool with qualifying BareMetalHosts")

			poolCRs, err := helper.ListResourcePoolCRs(HubAPIClient, tsparams.O2IMSNamespace)
			Expect(err).ToNot(HaveOccurred(), "Failed to list ResourcePool CRs")
			Expect(poolCRs).ToNot(BeEmpty(), "At least one ResourcePool CR should be present")

			var (
				selectedPool   = poolCRs[0]
				qualifiedHosts []*bmh.BmhBuilder
			)

			for _, poolCR := range poolCRs {
				hosts, listErr := helper.ListQualifiedBareMetalHostsForPool(HubAPIClient, poolCR.Name)
				Expect(listErr).ToNot(HaveOccurred(), "Failed to list BareMetalHosts for pool %s", poolCR.Name)

				if len(hosts) > 0 {
					selectedPool = poolCR
					qualifiedHosts = hosts

					break
				}
			}

			if len(qualifiedHosts) == 0 {
				Skip(inventoryBMHPrerequisite)
			}

			poolID, err := uuid.Parse(string(selectedPool.UID))
			Expect(err).ToNot(HaveOccurred(), "Failed to parse ResourcePool UID")

			By("listing resources for the selected resource pool")

			resources, err := inventoryClient.ListResources(poolID)
			Expect(err).ToNot(HaveOccurred(), "Failed to list resources from the API")
			Expect(resources).To(HaveLen(len(qualifiedHosts)), "API resource count should match qualifying BMH count")

			By("retrieving each resource by resourceId")

			for _, host := range qualifiedHosts {
				resource, found := helper.FindResourceByDescription(resources, host.Object.Name)
				Expect(found).To(BeTrue(), "Resource for BareMetalHost %s should be in list response", host.Object.Name)

				retrievedResource, err := inventoryClient.GetResource(poolID, resource.ResourceId)
				Expect(err).ToNot(HaveOccurred(), "Failed to get resource for BareMetalHost %s", host.Object.Name)
				Expect(retrievedResource.ResourcePoolId).To(Equal(poolID))
				Expect(retrievedResource.Description).To(Equal(host.Object.Name))
			}
		})

		// 00006 - Successfully list and retrieve resource types
		It("lists and retrieves resource types", reportxml.ID("00006"), func() {
			By("listing resource types from the API")

			resourceTypes, err := inventoryClient.ListResourceTypes()
			Expect(err).ToNot(HaveOccurred(), "Failed to list resource types from the API")
			if len(resourceTypes) == 0 {
				Skip(inventoryBMHPrerequisite)
			}

			By("retrieving each resource type by resourceTypeId")

			for _, resourceType := range resourceTypes {
				retrievedType, err := inventoryClient.GetResourceType(resourceType.ResourceTypeId)
				Expect(err).ToNot(HaveOccurred(), "Failed to get resource type %s", resourceType.ResourceTypeId)
				Expect(retrievedType.ResourceTypeId).To(Equal(resourceType.ResourceTypeId))
				Expect(retrievedType.Name).To(Equal(resourceType.Name))
			}
		})

		// 00007 - Successfully list and retrieve alarm dictionaries
		It("lists and retrieves alarm dictionaries", reportxml.ID("00007"), func() {
			By("listing alarm dictionaries from the API")

			alarmDictionaries, err := inventoryClient.ListAlarmDictionaries()
			Expect(err).ToNot(HaveOccurred(), "Failed to list alarm dictionaries from the API")
			if len(alarmDictionaries) == 0 {
				Skip(inventoryBMHPrerequisite)
			}

			By("retrieving each alarm dictionary by alarmDictionaryId")

			for _, alarmDictionary := range alarmDictionaries {
				retrievedDictionary, err := inventoryClient.GetAlarmDictionary(alarmDictionary.AlarmDictionaryId)
				Expect(err).ToNot(HaveOccurred(), "Failed to get alarm dictionary %s", alarmDictionary.AlarmDictionaryId)
				Expect(retrievedDictionary.AlarmDictionaryId).To(Equal(alarmDictionary.AlarmDictionaryId))
			}

			By("listing resource types with an associated alarm dictionary")

			resourceTypes, err := inventoryClient.ListResourceTypes()
			Expect(err).ToNot(HaveOccurred(), "Failed to list resource types from the API")

			resourceTypesWithDictionary := slices.DeleteFunc(slices.Clone(resourceTypes),
				func(resourceType oranapi.ResourceType) bool {
					return resourceType.AlarmDictionaryId == nil
				})
			if len(resourceTypesWithDictionary) == 0 {
				Skip(inventoryBMHPrerequisite)
			}

			By("retrieving the alarm dictionary for each associated resource type")

			for _, resourceType := range resourceTypesWithDictionary {
				alarmDictionary, err := inventoryClient.GetResourceTypeAlarmDictionary(resourceType.ResourceTypeId)
				Expect(err).ToNot(HaveOccurred(),
					"Failed to get alarm dictionary for resource type %s", resourceType.ResourceTypeId)
				Expect(alarmDictionary.AlarmDictionaryId).To(Equal(*resourceType.AlarmDictionaryId))
			}
		})

		// 00008 - Successfully list and retrieve deployment managers
		It("lists and retrieves deployment managers", reportxml.ID("00008"), func() {
			By("listing ManagedClusters on the hub cluster")

			managedClusters, err := helper.ListManagedClusters(HubAPIClient)
			Expect(err).ToNot(HaveOccurred(), "Failed to list ManagedClusters")
			Expect(managedClusters).ToNot(BeEmpty(), "At least one ManagedCluster should be present")

			By("listing deployment managers from the API")

			deploymentManagers, err := inventoryClient.ListDeploymentManagers()
			Expect(err).ToNot(HaveOccurred(), "Failed to list deployment managers from the API")
			Expect(deploymentManagers).To(HaveLen(len(managedClusters)),
				"API deployment manager count should match ManagedCluster count")

			By("retrieving each deployment manager by deploymentManagerId")

			for _, managedCluster := range managedClusters {
				matchingManager, found := helper.FindDeploymentManagerByName(deploymentManagers, managedCluster.Name)
				Expect(found).To(BeTrue(), "ManagedCluster %s should have a matching deployment manager",
					managedCluster.Name)

				retrievedManager, err := inventoryClient.GetDeploymentManager(matchingManager.DeploymentManagerId)
				Expect(err).ToNot(HaveOccurred(), "Failed to get deployment manager for %s", managedCluster.Name)
				Expect(retrievedManager.Name).To(Equal(managedCluster.Name))
			}
		})

		// 00009 - Successfully filter and select fields on list endpoints
		It("filters and selects fields on list endpoints", reportxml.ID("00009"), func() {
			By("listing deployment managers with fields=name")

			nameOnlyManagers, err := inventoryClient.ListDeploymentManagers(&oranapi.InventoryListParams{
				Fields: ptr.To("name"),
			})
			Expect(err).ToNot(HaveOccurred(), "Failed to list deployment managers with fields=name")
			Expect(nameOnlyManagers).ToNot(BeEmpty())

			allManagers, err := inventoryClient.ListDeploymentManagers()
			Expect(err).ToNot(HaveOccurred(), "Failed to list deployment managers")
			Expect(allManagers).ToNot(BeEmpty())

			for _, deploymentManager := range nameOnlyManagers {
				Expect(deploymentManager.Name).ToNot(BeEmpty())

				matchingManager, found := helper.FindDeploymentManagerByName(allManagers, deploymentManager.Name)
				Expect(found).To(BeTrue(), "fields=name response should include a known deployment manager")
				Expect(deploymentManager.DeploymentManagerId).To(Equal(matchingManager.DeploymentManagerId))
				Expect(deploymentManager.Extensions).To(BeNil())
			}

			By("listing deployment managers with exclude_fields=extensions")

			managersWithoutExtensions, err := inventoryClient.ListDeploymentManagers(&oranapi.InventoryListParams{
				ExcludeFields: ptr.To("extensions"),
			})
			Expect(err).ToNot(HaveOccurred(), "Failed to list deployment managers with exclude_fields=extensions")
			Expect(managersWithoutExtensions).ToNot(BeEmpty())

			for _, deploymentManager := range managersWithoutExtensions {
				Expect(deploymentManager.Extensions).To(BeNil())
			}

			By("filtering deployment managers by name")

			knownManager := allManagers[0]

			filteredManagers, err := inventoryClient.ListDeploymentManagers(&oranapi.InventoryListParams{
				Filter: filter.Equals("name", knownManager.Name),
			})
			Expect(err).ToNot(HaveOccurred(), "Failed to filter deployment managers by name")
			Expect(filteredManagers).To(HaveLen(1))
			Expect(filteredManagers[0].Name).To(Equal(knownManager.Name))

			By("filtering locations by globalLocationId")

			locationCRs, err := helper.ListLocationCRs(HubAPIClient, tsparams.O2IMSNamespace)
			Expect(err).ToNot(HaveOccurred(), "Failed to list Location CRs")
			Expect(locationCRs).ToNot(BeEmpty())

			knownLocation := locationCRs[0]

			filteredLocations, err := inventoryClient.ListLocations(&oranapi.InventoryListParams{
				Filter: filter.Equals("globalLocationId", knownLocation.Name),
			})
			Expect(err).ToNot(HaveOccurred(), "Failed to filter locations by globalLocationId")
			Expect(filteredLocations).To(HaveLen(1))
			Expect(filteredLocations[0].GlobalLocationId).To(Equal(knownLocation.Name))
		})

		// 00010 - Successfully create, retrieve, and delete a subscription
		It("creates, retrieves, and deletes a subscription", reportxml.ID("00010"), func() {
			By("creating a test subscription")

			consumerSubscriptionID := uuid.New()
			subscription, err := inventoryClient.CreateSubscription(oranapi.InventorySubscription{
				ConsumerSubscriptionId: &consumerSubscriptionID,
				Callback:               inventorySubscriberURL + "/" + consumerSubscriptionID.String(),
			})
			Expect(err).ToNot(HaveOccurred(), "Failed to create test subscription")
			Expect(subscription.SubscriptionId).ToNot(BeNil())
			Expect(*subscription.SubscriptionId).ToNot(Equal(uuid.Nil))
			Expect(subscription.ConsumerSubscriptionId.String()).To(Equal(consumerSubscriptionID.String()))
			Expect(subscription.Callback).To(Equal(inventorySubscriberURL + "/" + consumerSubscriptionID.String()))

			subscriptionID := *subscription.SubscriptionId

			By("listing subscriptions")

			subscriptions, err := inventoryClient.ListSubscriptions()
			Expect(err).ToNot(HaveOccurred(), "Failed to list subscriptions")

			containsSubscription := slices.ContainsFunc(subscriptions, func(item oranapi.InventorySubscription) bool {
				return item.SubscriptionId != nil && item.SubscriptionId.String() == subscriptionID.String()
			})
			Expect(containsSubscription).To(BeTrue(), "Created subscription should be in list response")

			By("retrieving the subscription")

			retrievedSubscription, err := inventoryClient.GetSubscription(subscriptionID)
			Expect(err).ToNot(HaveOccurred(), "Failed to retrieve test subscription")
			Expect(retrievedSubscription.ConsumerSubscriptionId.String()).To(Equal(consumerSubscriptionID.String()))

			By("deleting the subscription")

			err = inventoryClient.DeleteSubscription(subscriptionID)
			Expect(err).ToNot(HaveOccurred(), "Failed to delete test subscription")

			By("verifying the deleted subscription returns 404")

			_, err = inventoryClient.GetSubscription(subscriptionID)
			Expect(err).To(HaveOccurred(), "Deleted subscription should not be retrievable")

			apiErr := oranapi.AsAPIError(err)
			Expect(apiErr).ToNot(BeNil(), "Expected an API error for deleted subscription")
			Expect(apiErr.Status).To(Equal(404))
		})

		// 00011 - Successfully receive inventory change notification via subscription
		It("receives inventory change notification via subscription", reportxml.ID("00011"), func() {
			testLocationName := fmt.Sprintf("%s-inventory-%s", tsparams.TestName, uuid.NewString())

			By("creating a test subscription")

			consumerSubscriptionID := uuid.New()
			subscription, err := inventoryClient.CreateSubscription(oranapi.InventorySubscription{
				ConsumerSubscriptionId: &consumerSubscriptionID,
				Callback:               inventorySubscriberURL + "/" + consumerSubscriptionID.String(),
			})
			Expect(err).ToNot(HaveOccurred(), "Failed to create test subscription")
			Expect(subscription.SubscriptionId).ToNot(BeNil())

			subscriptionID := *subscription.SubscriptionId

			DeferCleanup(func() {
				By("deleting the test subscription")

				err := inventoryClient.DeleteSubscription(subscriptionID)
				Expect(err).ToNot(HaveOccurred(), "Failed to delete test subscription")

				By("deleting the test Location CR")

				err = helper.DeleteTestLocation(HubAPIClient, testLocationName)
				Expect(err).ToNot(HaveOccurred(), "Failed to delete test Location CR")
			})

			timeBeforeCreate := time.Now()

			By("creating a test Location CR")

			_, err = helper.CreateTestLocation(HubAPIClient, testLocationName, "inventory subscription test location")
			Expect(err).ToNot(HaveOccurred(), "Failed to create test Location CR")

			By("waiting for the Location to appear in the API")

			_, err = helper.WaitForLocationInAPI(inventoryClient, testLocationName, 2*time.Minute)
			Expect(err).ToNot(HaveOccurred(), "Failed to wait for test Location in API")

			By("waiting for the inventory change notification")

			err = subscriber.WaitForInventoryNotification(HubAPIClient, tsparams.SubscriberNamespace,
				subscriber.WithInventoryStart(timeBeforeCreate),
				subscriber.WithInventoryTimeout(2*time.Minute),
				subscriber.WithInventoryMatchFunc(func(notification *oranapi.InventoryChangeNotification) bool {
					if notification.NotificationEventType != oranapi.InventoryChangeNotificationTypeCreate {
						return false
					}

					if notification.PostObjectState == nil {
						return false
					}

					globalLocationID, ok := (*notification.PostObjectState)["globalLocationId"].(string)

					return ok && globalLocationID == testLocationName
				}),
			)
			Expect(err).ToNot(HaveOccurred(), "Failed to receive inventory change notification")
		})
	})
