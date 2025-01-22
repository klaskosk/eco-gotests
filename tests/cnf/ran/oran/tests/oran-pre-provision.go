package tests

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/oran"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-goinfra/pkg/secret"
	. "github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/helper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("ORAN Pre-provision Tests", Label(tsparams.LabelPreProvision), func() {
	// 77386 - Failed authentication with hardware manager
	It("fails to authenticate with hardware manager", reportxml.ID("77386"), func() {
		By("getting a valid Dell HardwareManager")
		hwmgr, err := helper.GetValidDellHwmgr(HubAPIClient)
		Expect(err).ToNot(HaveOccurred(), "Failed to get a valid Dell HardwareManager")

		By("getting the Dell AuthSecret")
		authSecret, err := secret.Pull(
			HubAPIClient, hwmgr.Definition.Spec.DellData.AuthSecret, tsparams.HardwareManagerNamespace)
		Expect(err).ToNot(HaveOccurred(), "Failed to get the HardwareManager AuthSecret")
		Expect(authSecret.Definition.Data).ToNot(BeNil(), "HardwareManager AuthSecret must have data")

		By("copying the secret and updating the password")
		authSecret.Definition.Name += "-test"
		authSecret.Definition.Data["password"] = []byte("d3JvbmdwYXNzd29yZA==") // wrongpassword

		authSecret, err = authSecret.Create()
		Expect(err).ToNot(HaveOccurred(), "Failed to create the new AuthSecret")

		By("copying the HardwareManager and updating the AuthSecret")
		hwmgr.Definition.Name += "-test"
		hwmgr.Definition.Spec.DellData.AuthSecret = authSecret.Definition.Name

		hwmgr, err = hwmgr.Create()
		Expect(err).ToNot(HaveOccurred(), "Failed to create the new HardwareManager")

		By("waiting for the authentication to fail")
		hwmgr, err = hwmgr.WaitForCondition(tsparams.HwmgrFailedAuthCondition, time.Minute)
		Expect(err).ToNot(HaveOccurred(), "Failed to wait for the HardwareManager to fail authentication")

		By("deleting the invalid HardwareManager")
		err = hwmgr.Delete()
		Expect(err).ToNot(HaveOccurred(), "Failed to delete the invalid HardwareManager")

		By("deleting the invalid AuthSecret")
		err = authSecret.Delete()
		Expect(err).ToNot(HaveOccurred(), "Failed to delete the invalid AuthSecret")
	})

	When("a ProvisioningRequest is applied", func() {
		var prBuilder *oran.ProvisioningRequestBuilder

		AfterEach(func() {
			if prBuilder != nil {
				err := prBuilder.DeleteAndWait(time.Minute)
				Expect(err).ToNot(HaveOccurred(), "Failed to delete the ProvisioningRequest")
			}
		})

		DescribeTable("ProvisionRequest pre-provision validations",
			func(templateVersion string, condition metav1.Condition) {
				prBuilder = createPR(templateVersion)
				prBuilder = checkPRStatus(prBuilder, condition)
			},
			// 77387 - Failed provisioning with nonexistent hardware profile
			Entry("fails to provision with nonexistent hardware profile",
				reportxml.ID("77387"), tsparams.TemplateNonexistentProfile, tsparams.PRHardwareProvisionFailedCondition),
			// 77388 - Failed provisioning with no hardware available
			Entry("fails to provision with no hardware available",
				reportxml.ID("77388"), tsparams.TemplateNoHardware, tsparams.PRHardwareProvisionFailedCondition),
			// 77389 - Failed provisioning with missing interface labels
			Entry("fails to provision with missing interface labels",
				reportxml.ID("77389"), tsparams.TemplateMissingLabels, tsparams.PRValidationFailedCondition),
			// 77390 - Failed provisioning with incorrect boot interface label
			Entry("fails to provision with incorrect boot interface label",
				reportxml.ID("77390"), tsparams.TemplateIncorrectLabel, tsparams.PRNodeConfigFailedCondition),
			// 77392 - Apply a ProvisioningRequest referencing an invalid ClusterTemplate
			Entry("fails to provision with invalid ClusterTemplate",
				reportxml.ID("77392"), tsparams.TemplateInvalid, tsparams.PRValidationFailedCondition),
			// 78245 - Missing schema while provisioning without hardware template
			Entry("fails to provision without a HardwareTemplate when required schema is missing",
				reportxml.ID("78245"), tsparams.TemplateMissingSchema, tsparams.PRValidationFailedCondition),
		)

		// 78246 - Successful provisioning without hardware template
		It("successfully generates ClusterInstance provisioning without HardwareTemplate", reportxml.ID("78246"), func() {
			By("creating a ProvisioningRequest")
			prBuilder = helper.NewNoTemplatePR(
				HubAPIClient, RANConfig.Spoke1Name, RANConfig.Spoke1Hostname, tsparams.TemplateNoHWTemplate)

			var err error
			prBuilder, err = prBuilder.Create()
			Expect(err).ToNot(HaveOccurred(), "Failed to create a ProvisioningRequest")

			By("waiting for its ClusterInstance to be processed")
			prBuilder, err = prBuilder.WaitForCondition(tsparams.PRCIProcesssedCondition, time.Minute)
			Expect(err).ToNot(HaveOccurred(), "Failed to wait for ClusterInstance to be processed")
		})
	})

})

func createPR(templateVersion string) *oran.ProvisioningRequestBuilder {
	By("creating a ProvisioningRequest")

	prBuilder := helper.NewProvisioningRequest(
		HubAPIClient, RANConfig.Spoke1Name, RANConfig.Spoke1Hostname, templateVersion)
	prBuilder, err := prBuilder.Create()
	Expect(err).ToNot(HaveOccurred(), "Failed to create a ProvisioningRequest")

	return prBuilder
}

func checkPRStatus(
	prBuilder *oran.ProvisioningRequestBuilder, condition metav1.Condition) *oran.ProvisioningRequestBuilder {
	By("checking the ProvisioningRequest status")

	prBuilder, err := prBuilder.WaitForCondition(condition, time.Minute)
	Expect(err).ToNot(HaveOccurred(), "Failed to verify the ProvisioningRequest status")

	return prBuilder
}
