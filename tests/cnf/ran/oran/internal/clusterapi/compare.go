package clusterapi

import (
	"fmt"
)

// appendMismatch appends a field mismatch error to errs when want and got differ.
func appendMismatch(errs []error, field string, want, got any) []error {
	if want != got {
		return append(errs, fmt.Errorf("%s: want %#v, got %#v", field, want, got))
	}

	return errs
}

// appendError appends err to errs when err is non-nil.
func appendError(errs []error, err error) []error {
	if err != nil {
		return append(errs, err)
	}

	return errs
}

// derefSlice returns the slice pointed to by s, treating a nil pointer to a slice as a nil slice.
func derefSlice[T any](s *[]T) []T {
	if s == nil {
		return nil
	}

	return *s
}

// ExtensionString returns the string form of extensions[key], or "" when missing.
func ExtensionString(extensions *map[string]any, key string) string {
	if extensions == nil {
		return ""
	}

	return extensionString(*extensions, key)
}

// extensionString returns the string form of extensions[key], or "" when missing.
func extensionString(extensions map[string]any, key string) string {
	if extensions == nil {
		return ""
	}

	value, ok := extensions[key]
	if !ok || value == nil {
		return ""
	}

	return fmt.Sprint(value)
}

// asStringKeyedMap converts common JSON-decoded map shapes to map[string]string.
func asStringKeyedMap(value any) (map[string]string, bool) {
	switch typed := value.(type) {
	case map[string]string:
		return typed, true
	case map[string]any:
		result := make(map[string]string, len(typed))
		for key, nested := range typed {
			result[key] = fmt.Sprint(nested)
		}

		return result, true
	default:
		return nil, false
	}
}
