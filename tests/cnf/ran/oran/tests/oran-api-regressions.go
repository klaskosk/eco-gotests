package tests

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api/fields"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api/filter"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/reportxml"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/auth"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
)

var _ = Describe("ORAN API Regressions", Label(tsparams.LabelPostProvision), func() {
	var (
		clusterClient   *oranapi.ClusterClient
		artifactsClient *oranapi.ArtifactsClient
	)

	BeforeEach(func() {
		By("creating the O2IMS API clients")

		clientBuilder, err := auth.NewClientBuilderForConfig(RANConfig)
		Expect(err).ToNot(HaveOccurred(), "Failed to create the O2IMS API client builder")

		clusterClient, err = clientBuilder.BuildCluster()
		Expect(err).ToNot(HaveOccurred(), "Failed to create the O2IMS cluster API client")

		artifactsClient, err = clientBuilder.BuildArtifacts()
		Expect(err).ToNot(HaveOccurred(), "Failed to create the O2IMS artifacts API client")
	})

	// OCPBUGS-105288 - Allow nested field selection and exclusion on map-typed extensions.
	It("supports nested fields on map-typed extensions", reportxml.ID("105288"), func() {
		By("listing the spoke NodeCluster to capture its original extensions")

		original, err := clusterClient.ListNodeClusters(
			oranapi.WithFilter(filter.Equals("name", RANConfig.Spoke1Name)))
		Expect(err).ToNot(HaveOccurred(), "Failed to list the spoke NodeCluster extensions")
		Expect(original).To(HaveLen(1), "Expected exactly one NodeCluster for spoke %s", RANConfig.Spoke1Name)
		Expect(original[0].Extensions).ToNot(BeNil(), "Spoke NodeCluster should include extensions")
		Expect(*original[0].Extensions).To(HaveKey("openshiftVersion"),
			"Original extensions should include openshiftVersion")

		By("listing the spoke NodeCluster with only extensions.openshiftVersion selected")

		selected, err := clusterClient.ListNodeClusters(
			oranapi.WithFilter(filter.Equals("name", RANConfig.Spoke1Name)),
			oranapi.WithFields(fields.Include("name", "extensions/openshiftVersion")))
		Expect(err).ToNot(HaveOccurred(), "Failed to select a nested NodeCluster extension field")
		Expect(selected).To(HaveLen(1), "Expected exactly one NodeCluster for spoke %s", RANConfig.Spoke1Name)
		Expect(selected[0].Name).To(Equal(RANConfig.Spoke1Name), "Selected NodeCluster should be the spoke")
		Expect(selected[0].Extensions).ToNot(BeNil(), "Selected NodeCluster should include extensions")
		Expect(*selected[0].Extensions).To(HaveKey("openshiftVersion"),
			"Selected extensions should include openshiftVersion")
		Expect(*selected[0].Extensions).To(HaveLen(1), "Only the requested nested extension should be returned")

		By("listing the spoke NodeCluster while excluding extensions.openshiftVersion")

		excluded, err := clusterClient.ListNodeClusters(
			oranapi.WithFilter(filter.Equals("name", RANConfig.Spoke1Name)),
			oranapi.WithFields(fields.Exclude("extensions/openshiftVersion")))
		Expect(err).ToNot(HaveOccurred(), "Failed to exclude a nested NodeCluster extension field")
		Expect(excluded).To(HaveLen(1), "Expected exactly one NodeCluster for spoke %s", RANConfig.Spoke1Name)
		Expect(excluded[0].Extensions).ToNot(BeNil(), "Other NodeCluster extensions should remain present")
		Expect(*excluded[0].Extensions).ToNot(HaveKey("openshiftVersion"),
			"Excluded extensions should not include openshiftVersion")

		expectedExtensions := maps.Clone(*original[0].Extensions)
		delete(expectedExtensions, "openshiftVersion")
		Expect(expectedExtensions).ToNot(BeEmpty(), "Original extensions should include a sibling extension")
		Expect(*excluded[0].Extensions).To(Equal(expectedExtensions),
			"Excluding openshiftVersion should preserve the spoke's other extensions")
	})

	// OCPBUGS-105287 - Preserve raw semicolon compound filters through query parsing and generated handlers.
	It("accepts compound filters with raw semicolons", reportxml.ID("105287"), func() {
		By("listing ManagedInfrastructureTemplates to choose a real name and version")

		templates, err := artifactsClient.ListManagedInfrastructureTemplates()
		Expect(err).ToNot(HaveOccurred(), "Failed to list ManagedInfrastructureTemplates")
		Expect(templates).ToNot(BeEmpty(), "At least one ManagedInfrastructureTemplate is required")

		type nameVersion struct {
			name    string
			version string
		}

		var chosen *oranapi.ManagedInfrastructureTemplate

		for index := range templates {
			if templates[index].Name != "" && templates[index].Version != "" {
				chosen = &templates[index]

				break
			}
		}

		Expect(chosen).ToNot(BeNil(), "At least one template with a name and version is required")

		chosenPair := nameVersion{name: chosen.Name, version: chosen.Version}

		var expectedPairs []nameVersion

		for _, template := range templates {
			if template.Name == chosenPair.name && template.Version == chosenPair.version {
				expectedPairs = append(expectedPairs, nameVersion{name: template.Name, version: template.Version})
			}
		}

		Expect(expectedPairs).ToNot(BeEmpty(), "The selected name/version pair should exist in the unfiltered list")

		// Escape values while preserving the literal semicolon between the filter terms.
		compoundFilter := fmt.Sprintf("(eq,name,%s);(eq,version,%s)",
			url.QueryEscape(chosen.Name), url.QueryEscape(chosen.Version))

		By("requesting the exact compound filter with a literal semicolon")

		response, err := artifactsClient.GetManagedInfrastructureTemplatesWithResponse(
			context.Background(), nil, func(_ context.Context, request *http.Request) error {
				request.URL.RawQuery = "filter=" + compoundFilter

				return nil
			})
		Expect(err).ToNot(HaveOccurred(), "Failed to send ManagedInfrastructureTemplates compound filter request")
		Expect(response.StatusCode()).To(Equal(http.StatusOK),
			"Compound filter should be accepted; response body: %s", response.Body)
		Expect(response.JSON200).ToNot(BeNil(), "Compound filter response should include a result list")

		actualPairs := make([]nameVersion, 0, len(*response.JSON200))
		for _, template := range *response.JSON200 {
			actualPairs = append(actualPairs, nameVersion{name: template.Name, version: template.Version})
		}

		Expect(actualPairs).To(ConsistOf(expectedPairs),
			"Compound filter should return the templates matching both selected name and version")
	})
})
