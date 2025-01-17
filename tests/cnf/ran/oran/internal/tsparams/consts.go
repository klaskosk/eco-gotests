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
)

const (
	// TemplateNonexistentProfile is the version associated with the nonexistent hardware profile test.
	TemplateNonexistentProfile = "2"
	// TemplateNoHardware is the version associated with the no hardware available test.
	TemplateNoHardware = "3"
	// TemplateMissingLabels is the version associated with the missing interface labels test.
	TemplateMissingLabels = "4"
	// TemplateIncorrectLabel is the version associated with the incorrect boot interface label test.
	TemplateIncorrectLabel = "5"
	// TemplateUpdateProfile is the version associated with the hardware profile update test.
	TemplateUpdateProfile = "6"
	// TemplateInvalid is the version associated with the invalid ClusterTemplate test.
	TemplateInvalid = "7"
	// TemplateUpdateDefaults is the version associated with the ClusterInstance defaults update test.
	TemplateUpdateDefaults = "8"
	// TemplateUpdateExisting is the version associated with the update existing PG manifest test.
	TemplateUpdateExisting = "9"
	// TemplateAddNew is the version associated with the add new manifest to existing PG test.
	TemplateAddNew = "10"
	// TemplateUpdateSchema is the version associated with the policyTemplateParameters schema update test.
	TemplateUpdateSchema = "11"
)
