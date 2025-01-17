package tsparams

import (
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/openshift-kni/k8sreporter"
	pluginv1alpha1 "github.com/openshift-kni/oran-hwmgr-plugin/api/hwmgr-plugin/v1alpha1"
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
)
