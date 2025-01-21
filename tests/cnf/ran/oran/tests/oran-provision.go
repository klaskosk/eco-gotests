package tests

import (
	"encoding/base64"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/configmap"
	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-goinfra/pkg/secret"
	"github.com/openshift-kni/eco-goinfra/pkg/siteconfig"
	. "github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/helper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
)

// ContinueOnError is deliberately left out of this Ordered container. If the invalid ProvisioningRequest does not
// become valid again, we cannot test provisioning with a valid ProvisioningRequest.
var _ = Describe("ORAN Provision Tests", Label(tsparams.LabelProvision), Ordered, func() {
	var prBuilder *oran.ProvisioningRequestBuilder

	AfterAll(func() {
		By("saving the spoke 1 admin kubeconfig")
		kubeconfigSecret, err := secret.Pull(HubAPIClient, RANConfig.Spoke1Name+"-admin-kubeconfig", RANConfig.Spoke1Name)
		Expect(err).ToNot(HaveOccurred(), "Failed to get the spoke 1 kubeconfig secret")

		kubeconfig, err := base64.StdEncoding.DecodeString(string(kubeconfigSecret.Definition.Data["kubeconfig"]))
		Expect(err).ToNot(HaveOccurred(), "Failed to decode the spoke 1 kubeconfig secret")

		err = os.WriteFile(RANConfig.Spoke2Kubeconfig, kubeconfig, 0644)
		Expect(err).ToNot(HaveOccurred(), "Failed to save the spoke 1 admin kubeconfig")
	})

	// 77393 - Apply a ProvisioningRequest with missing required input parameter
	It("recovers provisioning when invalid ProvisioningRequest is updated", reportxml.ID("77393"), func() {
		By("creating a ProvisioningRequest with nodeClusterName missing")
		prBuilder = oran.NewPRBuilder(
			HubAPIClient, RANConfig.Spoke1Name, tsparams.ClusterTemplateName, tsparams.TemplateValid).
			WithTemplateParameter("oCloudSiteId", RANConfig.Spoke1Name).
			WithTemplateParameter("policyTemplateParameters", map[string]any{}).
			WithTemplateParameter("clusterInstanceParameters", map[string]any{
				"clusterName": RANConfig.Spoke1Name,
				"nodes": []map[string]any{{
					"hostName": RANConfig.Spoke1Hostname,
				}},
			})

		var err error
		prBuilder, err = prBuilder.Create()
		Expect(err).ToNot(HaveOccurred(), "Failed to create a ProvisioningRequest with nodeClusterName missing")

		By("checking the ProvisioningRequest status for a failure")
		prBuilder, err = prBuilder.WaitForCondition(tsparams.PRValidationFailedCondition, time.Minute)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for the ProvisioningRequest to fail")

		By("updating the ProvisioningRequest to add nodeClusterName")
		prBuilder = prBuilder.WithTemplateParameter("nodeClusterName", RANConfig.Spoke1Name)
		prBuilder, err = prBuilder.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update the ProvisioningRequest to add nodeClusterName")

		By("waiting for ProvisioningRequest validation to succeed")
		prBuilder, err = prBuilder.WaitForCondition(tsparams.PRValidationSucceededCondition, time.Minute)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for ProvisioningRequest validation to succeed")
	})

	// 77394 - Apply a valid ProvisioningRequest
	It("successfully provisions with a valid ProvisioningRequest", reportxml.ID("77394"), func() {
		By("waiting for the ProvisioningRequest to apply configuration")
		var err error
		prBuilder, err = prBuilder.WaitForCondition(tsparams.PRConfigurationAppliedCondition, 2*time.Hour)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for the ProvisioningRequest to apply configuration")

		By("verifying all the policies are compliant")
		err = helper.WaitForPoliciesCompliant(HubAPIClient, RANConfig.Spoke1Name, time.Minute)
		Expect(err).ToNot(HaveOccurred(), "Failed to verify all spoke 1 policies are compliant")

		By("verifying a NodePool was created")
		nodePool, err := oran.PullNodePool(HubAPIClient, RANConfig.Spoke1Name, tsparams.HardwareManagerNamespace)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull NodePool for spoke 1")
		Expect(nodePool.Object).ToNot(BeNil(), "Failed to get NodePool object for spoke 1")

		By("verifying spoke 1 pull-secret was created")
		pullSecret, err := secret.Pull(HubAPIClient, "pull-secret", RANConfig.Spoke1Name)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull pull-secret for spoke 1")
		Expect(pullSecret.Object).ToNot(BeNil(), "Failed to get pull-secret object for spoke 1")

		By("verifying spoke 1 extra-manifests was created")
		extraManifests, err := configmap.Pull(HubAPIClient, tsparams.ExtraManifestsName, RANConfig.Spoke1Name)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull extra-manifests ConfigMap for spoke 1")
		Expect(extraManifests.Object).ToNot(BeNil(), "Failed to get extra-manifests ConfigMap object for spoke 1")

		By("verifying spoke 1 policy ConfigMap was created")
		pgMap, err := configmap.Pull(
			HubAPIClient, RANConfig.Spoke1Name+"-pg", "ztp-"+getClusterTemplateNamespace(HubAPIClient))
		Expect(err).ToNot(HaveOccurred(), "Failed to pull policy ConfigMap for spoke 1")
		Expect(pgMap.Object).ToNot(BeNil(), "Failed to get policy ConfigMap object for spoke 1")

		if RANConfig.Spoke1BMC == nil {
			By("skipping check since spoke 1 BMC details are not provided")

			return
		}

		By("verifying the ClusterInstance has the correct BMC details")
		clusterInstance, err := siteconfig.PullClusterInstance(HubAPIClient, RANConfig.Spoke1Name, RANConfig.Spoke1Name)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull ClusterInstance for spoke 1")

		clusterInstanceNode := clusterInstance.Definition.Spec.Nodes[0]
		Expect(clusterInstanceNode.BmcAddress).
			To(ContainSubstring(RANConfig.BMCHosts[0]), "ClusterInstance has incorrect BMC address")

		bmcSecret, err := secret.Pull(HubAPIClient, clusterInstanceNode.BmcCredentialsName.Name, RANConfig.Spoke1Name)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull spoke 1 BMC secret")

		bmcUsername, err := base64.StdEncoding.DecodeString(string(bmcSecret.Definition.Data["username"]))
		Expect(err).ToNot(HaveOccurred(), "Failed to decode spoke 1 BMC username")
		Expect(bmcUsername).To(Equal(RANConfig.BMCUsername), "ClusterInstance has incorrect BMC username")

		bmcPassword, err := base64.StdEncoding.DecodeString(string(bmcSecret.Definition.Data["password"]))
		Expect(err).ToNot(HaveOccurred(), "Failed to decode spoke 1 BMC password")
		Expect(bmcPassword).To(Equal(RANConfig.BMCPassword), "ClusterInstance has incorrect BMC password")
	})
})

func getClusterTemplateNamespace(client *clients.Settings) string {
	clusterTemplates, err := oran.ListClusterTemplates(client)
	Expect(err).ToNot(HaveOccurred(), "Failed to list ClusterTemplates")

	expectedName := tsparams.ClusterTemplateName + "." + tsparams.TemplateValid
	for _, clusterTemplate := range clusterTemplates {
		if clusterTemplate.Definition.Name == expectedName {
			return clusterTemplate.Definition.Namespace
		}
	}

	return ""
}
