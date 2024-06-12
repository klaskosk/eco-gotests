package ranparam

import "github.com/openshift-kni/eco-gotests/tests/cnf/internal/cnfparams"

var (
	// Labels represents the range of labels that can be used for test cases selection.
	Labels = []string{cnfparams.Label, Label}

	// Spoke1Name is the name of the first spoke cluster.
	Spoke1Name string
	// Spoke2Name is the name of the second spoke cluster.
	Spoke2Name string

	// AcmVersion is the version of the ACM operator.
	AcmVersion string
	// TalmVersion is the version of the TALM operator.
	TalmVersion string
	// ZtpVersion is the version of the ZTP from ArgoCD.
	ZtpVersion string
)
