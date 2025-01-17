package tests

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-goinfra/pkg/secret"
	. "github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/helper"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
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

	// 77387 - Failed provisioning with nonexistent hardware profile
	It("fails to provision with nonexistent hardware profile", reportxml.ID("77387"), func() {
		By("creating a ProvisioningRequest")
		By("checking the ProvisioningRequest status")
	})

	// 77388 - Failed provisioning with no hardware available
	It("fails to provision with no hardware available", reportxml.ID("77388"), func() {
		By("creating a ProvisioningRequest")
		By("checking the ProvisioningRequest status")
	})

	// 77389 - Failed provisioning with missing interface labels
	It("fails to provision with missing interface labels", reportxml.ID("77389"), func() {
		By("creating a ProvisioningRequest")
		By("checking the ProvisioningRequest status")
	})

	// 77390 - Failed provisioning with incorrect boot interface label
	It("fails to provision with incorrect boot interface label", reportxml.ID("77390"), func() {
		By("creating a ProvisioningRequest")
		By("checking the ProvisioningRequest status")
	})

	// 77392 - Apply a ProvisioningRequest referencing an invalid ClusterTemplate
	It("fails to provision with invalid ClusterTemplate", reportxml.ID("77392"), func() {

	})
})
