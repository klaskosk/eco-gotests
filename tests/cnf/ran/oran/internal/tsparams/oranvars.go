package tsparams

import (
	"github.com/openshift-kni/eco-gotests/tests/cnf/ran/internal/ranparam"
	"github.com/openshift-kni/k8sreporter"
)

var (
	// Labels is the labels applied to all test cases in the suite.
	Labels = append(ranparam.Labels, LabelSuite)

	// ReporterHubNamespacesToDump tells the reporter which namespaces on the hub to collect pod logs from.
	ReporterHubNamespacesToDump = map[string]string{}

	// ReporterSpokeNamespacesToDump tells the reporter which namespaces on the spokes to collect pod logs from.
	ReporterSpokeNamespacesToDump = map[string]string{}

	// ReporterHubCRsToDump is the CRs the reporter should dump on the hub.
	ReporterHubCRsToDump = []k8sreporter.CRData{}

	// ReporterSpokeCRsToDump is the CRs the reporter should dump on the spokes.
	ReporterSpokeCRsToDump = []k8sreporter.CRData{}
)
