package tests

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/configmap"
	"github.com/openshift-kni/eco-goinfra/pkg/ocm"
	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-goinfra/pkg/siteconfig"
	. "github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/helper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"github.com/openshift-kni/eco-gotests/tests/internal/cluster"
	provisioningv1alpha1 "github.com/openshift-kni/oran-o2ims/api/provisioning/v1alpha1"
)

var _ = Describe("ORAN Post-provision Tests", Label(tsparams.LabelPostProvision), func() {
	var (
		prBuilder      *oran.ProvisioningRequestBuilder
		originalPRSpec *provisioningv1alpha1.ProvisioningRequestSpec
	)

	BeforeEach(func() {
		By("saving the original ProvisioningRequest spec")
		var err error
		prBuilder, err = oran.PullPR(HubAPIClient, RANConfig.Spoke1Name)
		Expect(err).ToNot(HaveOccurred(), "Failed to pull spoke 1 ProvisioningRequest")

		copiedSpec := prBuilder.Definition.Spec
		originalPRSpec = &copiedSpec
	})

	AfterEach(func() {
		// If saving the original spec failed, skip restoring it.
		if originalPRSpec == nil {
			return
		}

		By("restoring the original ProvisioningRequest spec")
		prBuilder.Definition.Spec = *originalPRSpec
		_, err := prBuilder.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to restore spoke 1 ProvisioningRequest")

		By("waiting for spoke 1 to recover")
		err = cluster.WaitForRecover(Spoke1APIClient, []string{}, 45*time.Minute)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for spoke 1 to recover")
	})

	// 77373 - Successful update to ProvisioningRequest clusterInstanceParameters
	It("successfully updates clusterInstanceParameters", reportxml.ID("77373"), func() {
		By("verifying the label does not exist already")
		verifyNoLabel()

		By("updating the extraLabels in clusterInstanceParameters")
		templateParameters, err := prBuilder.GetTemplateParameters()
		Expect(err).ToNot(HaveOccurred(), "Failed to get spoke 1 TemplateParameters")
		Expect(tsparams.ClusterInstanceParamsKey).
			To(BeKeyOf(templateParameters), "Spoke 1 TemplateParameters is missing clusterInstanceParameters")

		clusterInstanceParams, ok := templateParameters[tsparams.ClusterInstanceParamsKey].(map[string]any)
		Expect(ok).To(BeTrue(), "Spoke 1 clusterInstanceParameters is not a map[string]any")

		clusterInstanceParams["extraLabels"] = map[string]any{"ManagedCluster": map[string]string{tsparams.TestName: ""}}
		err = prBuilder.WithTemplateParameters(templateParameters)
		Expect(err).ToNot(HaveOccurred(), "Failed to set spoke 1 TemplateParameters")

		_, err = prBuilder.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update spoke 1 ProvisioningRequest")

		waitForLabels()
	})

	// 77374 - Successful update to ProvisioningRequest policyTemplateParameters
	It("successfully updates policyTemplateParameters", reportxml.ID("77374"), func() {
		By("verifying the test ConfigMap exists and has the original value")
		verifyCM(tsparams.TestName, tsparams.TestOriginalValue)

		By("updating the policyTemplateParameters")
		prBuilder = prBuilder.WithTemplateParameter(tsparams.PolicyTemplateParamsKey, map[string]string{
			tsparams.TestName: tsparams.TestNewValue,
		})

		prBuilder, err := prBuilder.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update spoke 1 ProvisioningRequest")

		waitForPolicies(prBuilder)

		By("verifying the test ConfigMap has the new value")
		verifyCM(tsparams.TestName, tsparams.TestNewValue)
	})

	// 77375 - Successful update to ClusterInstance defaults ConfigMap
	It("successfully updates ClusterInstance defaults", reportxml.ID("77375"), func() {
		By("verifying the label does not exist already")
		verifyNoLabel()

		By("updating the ProvisioningRequest TemplateVersion")
		prBuilder.Definition.Spec.TemplateVersion = tsparams.TemplateUpdateDefaults
		_, err := prBuilder.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update spoke 1 ProvisioningRequest")

		waitForLabels()
	})

	// 77376 - Successful update of existing PG manifest
	It("successfully updates existing PG manifest", reportxml.ID("77376"), func() {
		By("verifying the test ConfigMap exists and has the original value")
		verifyCM(tsparams.TestName, tsparams.TestOriginalValue)

		By("updating the ProvisioningRequest TemplateVersion")
		prBuilder.Definition.Spec.TemplateVersion = tsparams.TemplateUpdateExisting
		_, err := prBuilder.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update spoke 1 ProvisioningRequest")

		waitForPolicies(prBuilder)

		By("verifying the test ConfigMap has the new value")
		verifyCM(tsparams.TestName, tsparams.TestNewValue)
	})

	// 77377 - Successful addition of new manifest to existing PG
	It("successfully adds new manifest to existing PG", reportxml.ID("77377"), func() {
		By("verifying the test ConfigMap exists and has the original value")
		verifyCM(tsparams.TestName, tsparams.TestOriginalValue)

		By("verifying the second test ConfigMap does not exist")
		_, err := configmap.Pull(Spoke1APIClient, tsparams.TestName2, tsparams.TestName)
		Expect(err).To(HaveOccurred(), "Second test ConfigMap already exists on spoke 1")

		By("updating the ProvisioningRequest TemplateVersion")
		prBuilder.Definition.Spec.TemplateVersion = tsparams.TemplateAddNew
		_, err = prBuilder.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update spoke 1 ProvisioningRequest")

		waitForPolicies(prBuilder)

		By("verifying the test ConfigMap has the original value")
		verifyCM(tsparams.TestName, tsparams.TestOriginalValue)

		By("verifying the second test ConfigMap exists and has the original value")
		verifyCM(tsparams.TestName2, tsparams.TestOriginalValue)
	})

	// 77378 - Successful update of ClusterTemplate policyTemplateParameters schema
	It("successfully updates schema of policyTemplateParameters", reportxml.ID("77378"), func() {
		By("verifying the test ConfigMap exists and has the original value")
		verifyCM(tsparams.TestName, tsparams.TestOriginalValue)

		By("verifying the second test ConfigMap does not exist")
		_, err := configmap.Pull(Spoke1APIClient, tsparams.TestName2, tsparams.TestName)
		Expect(err).To(HaveOccurred(), "Second test ConfigMap already exists on spoke 1")

		By("updating the ProvisioningRequest TemplateVersion")
		prBuilder.Definition.Spec.TemplateVersion = tsparams.TemplateUpdateSchema
		_, err = prBuilder.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update spoke 1 ProvisioningRequest")

		waitForPolicies(prBuilder)

		By("verifying the test ConfigMap has the original value")
		verifyCM(tsparams.TestName, tsparams.TestOriginalValue)

		By("verifying the second test ConfigMap has the new value")
		verifyCM(tsparams.TestName2, tsparams.TestNewValue)
	})

	// 77379 - Failed update to ProvisioningRequest and successful rollback
	It("successfully rolls back failed ProvisioningRequest update", reportxml.ID("77379"), func() {
		By("verifying ProvisioningRequest is valid to start")
		prBuilder, err := prBuilder.WaitForCondition(tsparams.PRConfigurationAppliedCondition, time.Minute)
		Expect(err).ToNot(HaveOccurred(), "Failed to verify spoke 1 ProvisioningRequest has ConfigurationApplied")

		By("updating the policyTemplateParameters")
		prBuilder = prBuilder.WithTemplateParameter(tsparams.PolicyTemplateParamsKey, map[string]string{
			tsparams.HugePagesSizeKey: "2G",
		})

		prBuilder, err = prBuilder.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update spoke 1 ProvisioningRequest")

		By("waiting for policy to go NonCompliant")
		err = helper.WaitForNoncompliantImmutable(HubAPIClient, RANConfig.Spoke1Name, time.Minute)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for a spoke 1 policy to go NonCompliant due to immutable field")

		By("fixing the policyTemplateParameters")
		prBuilder = prBuilder.WithTemplateParameter(tsparams.PolicyTemplateParamsKey, map[string]string{})
		prBuilder, err = prBuilder.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update spoke 1 ProvisioningRequest")

		waitForPolicies(prBuilder)
	})

	// 77391 - Successful update of hardware profile
	It("successfully updates hardware profile", reportxml.ID("77391"), func() {
		By("verifying spoke 1 is powered on")
		powerState, err := BMCClient.SystemPowerState()
		Expect(err).ToNot(HaveOccurred(), "Failed to get system power state from spoke 1 BMC")
		Expect(powerState).To(Equal("On"), "Spoke 1 is not powered on")

		By("updating ProvisioningRequest TemplateVersion")
		prBuilder.Definition.Spec.TemplateVersion = tsparams.TemplateUpdateProfile
		_, err = prBuilder.Update()
		Expect(err).ToNot(HaveOccurred(), "Failed to update spoke 1 ProvisioningRequest")

		By("waiting for spoke 1 to be powered off")
		err = helper.WaitForPoweredOff(BMCClient, 5*time.Minute)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for spoke 1 to power off")
	})
})

// verifyCM verifies that the test ConfigMap name has value for the test key.
func verifyCM(name, value string) {
	testCM, err := configmap.Pull(Spoke1APIClient, name, tsparams.TestName)
	Expect(err).ToNot(HaveOccurred(), "Failed to pull test ConfigMap %s from spoke 1", name)
	Expect(tsparams.TestName).
		To(BeKeyOf(testCM.Definition.Data), "Test ConfigMap %s on spoke 1 does not have test key", name)
	Expect(testCM.Definition.Data[tsparams.TestName]).
		To(Equal(value), "Test ConfigMap %s on spoke 1 does not have value %s", name, value)
}

// verifyNoLabel verifies that the ManagedCluster does not already have the test label.
func verifyNoLabel() {
	mcl, err := ocm.PullManagedCluster(HubAPIClient, RANConfig.Spoke1Name)
	Expect(err).ToNot(HaveOccurred(), "Failed to pull spoke 1 ManagedCluster")

	_, hasLabel := mcl.Definition.Labels[tsparams.TestName]
	Expect(hasLabel).To(BeFalse(), "Spoke 1 managed cluster already has label")
}

// waitForPolicies waits first for the policies to compliant then for prBuilder to have the ConfigurationApplied
// condition.
func waitForPolicies(prBuilder *oran.ProvisioningRequestBuilder) {
	By("waiting for policies to be compliant")

	err := helper.WaitForPoliciesCompliant(HubAPIClient, RANConfig.Spoke1Name, time.Minute)
	Expect(err).ToNot(HaveOccurred(), "Failed to wait for spoke 1 policies to be compliant")

	By("verifying the ProvisioningRequest status is updated")

	_, err = prBuilder.WaitForCondition(tsparams.PRConfigurationAppliedCondition, time.Minute)
	Expect(err).ToNot(HaveOccurred(), "Failed to wait for spoke 1 ProvisioningRequest to have ConfigurationApplied")
}

// waitForLabels waits for the test label to appear on the ClusterInstance then on the ManagedCluster.
func waitForLabels() {
	By("waiting for ClusterInstance to have label")

	clusterInstance, err := siteconfig.PullClusterInstance(HubAPIClient, RANConfig.Spoke1Name, RANConfig.Spoke1Name)
	Expect(err).ToNot(HaveOccurred(), "Failed to pull spoke 1 ClusterInstance")

	err = helper.WaitForCIExtraLabel(clusterInstance, tsparams.TestName, time.Minute)
	Expect(err).ToNot(HaveOccurred(), "Failed to wait for spoke 1 ClusterInstance to have the extraLabel")

	By("waiting for ManagedCluster to have label")

	mcl, err := ocm.PullManagedCluster(HubAPIClient, RANConfig.Spoke1Name)
	Expect(err).ToNot(HaveOccurred(), "Failed to pull spoke 1 ManagedCluster")

	err = helper.WaitForMCLLabel(mcl, tsparams.TestName, time.Minute)
	Expect(err).ToNot(HaveOccurred(), "Failed to wait for spoke 1 ManagedCluster to have the label")
}
