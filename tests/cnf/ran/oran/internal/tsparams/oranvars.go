package tsparams

import (
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/openshift-kni/k8sreporter"
	pluginv1alpha1 "github.com/openshift-kni/oran-hwmgr-plugin/api/hwmgr-plugin/v1alpha1"
	provisioningv1alpha1 "github.com/openshift-kni/oran-o2ims/api/provisioning/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// Labels is the labels applied to all test cases in the suite.
	Labels = append(ranparam.Labels, LabelSuite)

	// ReporterHubNamespacesToDump tells the reporter which namespaces on the hub to collect pod logs from.
	ReporterHubNamespacesToDump = map[string]string{}

	// ReporterSpokeNamespacesToDump tells the reporter which namespaces on the spoke to collect pod logs from.
	ReporterSpokeNamespacesToDump = map[string]string{}

	// ReporterHubCRsToDump is the CRs the reporter should dump on the hub.
	ReporterHubCRsToDump = []k8sreporter.CRData{}

	// ReporterSpokeCRsToDump is the CRs the reporter should dump on the spoke.
	ReporterSpokeCRsToDump = []k8sreporter.CRData{}
)

var (
	// HwmgrFailedAuthCondition is the condition to match for when the HardwareManager fails to authenticate with
	// the DTIAS.
	HwmgrFailedAuthCondition = metav1.Condition{
		Type:    string(pluginv1alpha1.ConditionTypes.Validation),
		Reason:  string(pluginv1alpha1.ConditionReasons.Failed),
		Status:  metav1.ConditionFalse,
		Message: "401",
	}

	// PRHardwareProvisionFailedCondition is the ProvisioningRequest condition where hardware provisioning failed.
	PRHardwareProvisionFailedCondition = metav1.Condition{
		Type:   string(provisioningv1alpha1.PRconditionTypes.HardwareProvisioned),
		Reason: string(provisioningv1alpha1.CRconditionReasons.Failed),
		Status: metav1.ConditionFalse,
	}
	// PRValidationFailedCondition is the ProvisioningRequest condition where ProvisioningRequest validation failed.
	PRValidationFailedCondition = metav1.Condition{
		Type:   string(provisioningv1alpha1.PRconditionTypes.Validated),
		Reason: string(provisioningv1alpha1.CRconditionReasons.Failed),
		Status: metav1.ConditionFalse,
	}
	// PRValidationSucceededCondition is the ProvisioningRequest condition where ProvisioningRequest validation
	// succeeded.
	PRValidationSucceededCondition = metav1.Condition{
		Type:   string(provisioningv1alpha1.PRconditionTypes.Validated),
		Reason: string(provisioningv1alpha1.CRconditionReasons.Completed),
		Status: metav1.ConditionTrue,
	}
	// PRNodeConfigFailedCondition is the ProvisioningRequest condition where applying the node configuration
	// failed.
	PRNodeConfigFailedCondition = metav1.Condition{
		Type:   string(provisioningv1alpha1.PRconditionTypes.HardwareNodeConfigApplied),
		Reason: string(provisioningv1alpha1.CRconditionReasons.NotApplied),
		Status: metav1.ConditionFalse,
	}
	// PRConfigurationAppliedCondition is the ProvisioningRequest condition where applying day2 configuration
	// succeeds.
	PRConfigurationAppliedCondition = metav1.Condition{
		Type:   string(provisioningv1alpha1.PRconditionTypes.ConfigurationApplied),
		Reason: string(provisioningv1alpha1.CRconditionReasons.Completed),
		Status: metav1.ConditionTrue,
	}
	// PRCIProcesssedCondition is the ProvisioningRequest condition where the ClusterInstance has successfully been
	// processed.
	PRCIProcesssedCondition = metav1.Condition{
		Type:   string(provisioningv1alpha1.PRconditionTypes.ClusterInstanceProcessed),
		Reason: string(provisioningv1alpha1.CRconditionReasons.Completed),
		Status: metav1.ConditionTrue,
	}
)
