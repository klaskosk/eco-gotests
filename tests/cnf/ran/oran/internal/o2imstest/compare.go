package o2imstest

import (
	"fmt"
)

// AppendMismatch appends a field mismatch error to errs when want and got differ.
func AppendMismatch(errs []error, field string, want, got any) []error {
	if want != got {
		return append(errs, fmt.Errorf("%s: want %#v, got %#v", field, want, got))
	}

	return errs
}

// AppendError appends err to errs when err is non-nil.
func AppendError(errs []error, err error) []error {
	if err != nil {
		return append(errs, err)
	}

	return errs
}
