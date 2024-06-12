package tsparams

import (
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/openshift-kni/k8sreporter"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

// ArgoCdGitDetails is the details for a single app in ArgoCD.
type ArgoCdGitDetails struct {
	Repo   string
	Branch string
	Path   string
}

var (
	// Labels represents the range of labels that can be used for test cases selection.
	Labels = append(ranparam.Labels, LabelSuite)

	// ReporterNamespacesToDump tells to the reporter from where to collect logs.
	ReporterNamespacesToDump = map[string]string{
		TestNamespace: "",
	}
	// ReporterCRDsToDump tells to the reporter what CRs to dump.
	ReporterCRDsToDump = []k8sreporter.CRData{
		{Cr: &corev1.PodList{}},
	}

	// ArgoCdApps is the slice of the ArgoCd app names defined in this package.
	ArgoCdApps = []string{
		ArgoCdClustersAppName,
		ArgoCdPoliciesAppName,
	}
	// ArgoCdAppDetails contains more details for each of the ArgoCdApps.
	ArgoCdAppDetails = map[string]ArgoCdGitDetails{}

	// ImageRegistryPolicies is a slice of all the policies the image registry test creates.
	ImageRegistryPolicies = []string{
		"image-registry-policy-sc",
		"image-registry-policy-pvc",
		"image-registry-policy-pv",
		"image-registry-policy-config",
	}

	// NonManagementNamespaces is all the namespaces to consider for testing and not cluster management.
	NonManagementNamespaces = sets.NewString(
		TestNamespace, "cnf-ran-gotests-priv", "vran-acceleration-operators", "amq-router")

	// InvalidManagedPoliciesCondition is the CGU condition for where there are invalid managed policies.
	InvalidManagedPoliciesCondition = metav1.Condition{
		Type:    "Validated",
		Status:  metav1.ConditionFalse,
		Reason:  "NotAllManagedPoliciesExist",
		Message: "Invalid managed policies",
	}
)
