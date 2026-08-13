package o2imstest

import (
	"errors"
	"fmt"

	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
)

// VerifyAPIVersions checks that versions matches the expected first API version and URI prefix.
func VerifyAPIVersions(versions oranapi.APIVersions, expectedVersion, expectedURIPrefix string) error {
	var errs []error

	if versions.ApiVersions == nil || len(*versions.ApiVersions) == 0 || (*versions.ApiVersions)[0].Version == nil {
		errs = append(errs, fmt.Errorf("apiVersions[0].version: want non-nil, got nil"))
	} else {
		errs = AppendMismatch(errs, "apiVersions[0].version", expectedVersion, *(*versions.ApiVersions)[0].Version)
	}

	if versions.UriPrefix == nil {
		errs = append(errs, fmt.Errorf("uriPrefix: want non-nil, got nil"))
	} else {
		errs = AppendMismatch(errs, "uriPrefix", expectedURIPrefix, *versions.UriPrefix)
	}

	return errors.Join(errs...)
}
