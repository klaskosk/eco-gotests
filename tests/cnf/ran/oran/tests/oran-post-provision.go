package tests

import (
	. "github.com/onsi/ginkgo/v2"
	// . "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/reportxml"
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
)

var _ = Describe("ORAN Post-provision Tests", Label(tsparams.LabelPostProvision), func() {
	// 77373 - Successful update to ProvisioningRequest clusterInstanceParameters
	It("successfully updates clusterInstanceParameters", reportxml.ID("77373"), func() {

	})

	// 77374 - Successful update to ProvisioningRequest policyTemplateParameters
	It("successfully updates policyTemplateParameters", reportxml.ID("77374"), func() {

	})

	// 77375 - Successful update to ClusterInstance defaults ConfigMap
	It("successfully updates ClusterInstance defaults", reportxml.ID("77375"), func() {

	})

	// 77376 - Successful update of existing PG manifest
	It("successfully updates existing PG manifest", reportxml.ID("77376"), func() {

	})

	// 77377 - Successful addition of new manifest to existing PG
	It("successfully adds new manifest to existing PG", reportxml.ID("77377"), func() {

	})

	// 77378 - Successful update of ClusterTemplate policyTemplateParameters schema
	It("successfully updates schema of policyTemplateParameters", reportxml.ID("77378"), func() {

	})

	// 77379 - Failed update to ProvisioningRequest and successful rollback
	It("successfully rolls back failed ProvisioningRequest update", reportxml.ID("77379"), func() {

	})

	// 77391 - Successful update of hardware profile
	It("successfully updates hardware profile", reportxml.ID("77391"), func() {

	})
})
