package tsparams

const (
	// LabelSuite is the label applied to all cases in the oran suite.
	LabelSuite = "oran"
	// LabelPreProvision is the label applied to just the pre-provision test cases.
	LabelPreProvision = "pre-provision"
	// LabelProvision is the label applied to just the provision test cases.
	LabelProvision = "provision"
	// LabelPostProvision is the label applied to just the post-provision test cases.
	LabelPostProvision = "post-provision"
)

const (
	// ClusterTemplateName is the name without the version of the ClusterTemplate used in the ORAN tests.
	ClusterTemplateName = "sno-ran-du"
	// HardwareManagerNamespace is the namespace that HardwareManagers and their secrets use.
	HardwareManagerNamespace = "oran-hwmgr-plugin"
	// ExtraManifestsName is the name of the generated extra manifests ConfigMap in the cluster Namespace.
	ExtraManifestsName = "sno-ran-du-extra-manifests-1"
)

const (
	// TemplateValid is the valid ClusterTemplate used for the provision tests.
	TemplateValid = "v1"
	// TemplateNonexistentProfile is the version associated with the nonexistent hardware profile test.
	TemplateNonexistentProfile = "v2"
	// TemplateNoHardware is the version associated with the no hardware available test.
	TemplateNoHardware = "v3"
	// TemplateMissingLabels is the version associated with the missing interface labels test.
	TemplateMissingLabels = "v4"
	// TemplateIncorrectLabel is the version associated with the incorrect boot interface label test.
	TemplateIncorrectLabel = "v5"
	// TemplateUpdateProfile is the version associated with the hardware profile update test.
	TemplateUpdateProfile = "v6"
	// TemplateInvalid is the version associated with the invalid ClusterTemplate test.
	TemplateInvalid = "v7"
	// TemplateUpdateDefaults is the version associated with the ClusterInstance defaults update test.
	TemplateUpdateDefaults = "v8"
	// TemplateUpdateExisting is the version associated with the update existing PG manifest test.
	TemplateUpdateExisting = "v9"
	// TemplateAddNew is the version associated with the add new manifest to existing PG test.
	TemplateAddNew = "v10"
	// TemplateUpdateSchema is the version associated with the policyTemplateParameters schema update test.
	TemplateUpdateSchema = "v11"
)
