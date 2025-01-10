package tests

import (
	. "github.com/onsi/ginkgo/v2"
	// . "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
)

var _ = Describe("ORAN Pre-provision Tests", Label(tsparams.LabelPreProvision), func() {
	// 77386 - Failed authentication with hardware manager
	It("fails to authenticate with hardware manager", reportxml.ID("77386"), func() {

	})

	// 77387 - Failed provisioning with nonexistent hardware profile
	It("fails to provision with nonexistent hardware profile", reportxml.ID("77387"), func() {

	})

	// 77388 - Failed provisioning with no hardware available
	It("fails to provision with no hardware available", reportxml.ID("77388"), func() {

	})

	// 77389 - Failed provisioning with missing interface labels
	It("fails to provision with missing interface labels", reportxml.ID("77389"), func() {

	})

	// 77390 - Failed provisioning with incorrect boot interface label
	It("fails to provision with incorrect boot interface label", reportxml.ID("77390"), func() {

	})

	// 77391 - Successful update of hardware profile
	It("successfully updates hardware profile", reportxml.ID("77391"), func() {

	})

	// 77392 - Apply a ProvisioningRequest referencing an invalid ClusterTemplate
	It("fails to provision with invalid ClusterTemplate", reportxml.ID("77392"), func() {

	})
})
