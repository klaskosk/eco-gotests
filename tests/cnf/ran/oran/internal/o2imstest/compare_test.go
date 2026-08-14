package o2imstest

import (
	"fmt"
	"testing"

	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"github.com/stretchr/testify/assert"
)

func TestAsStringKeyedMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     any
		converted bool
		expected  map[string]string
	}{
		{
			name:      "map string any",
			value:     map[string]any{"cores": "96", "architecture": "x86_64"},
			converted: true,
			expected:  map[string]string{"cores": "96", "architecture": "x86_64"},
		},
		{
			name:      "map string string",
			value:     map[string]string{"GiB": "256"},
			converted: true,
			expected:  map[string]string{"GiB": "256"},
		},
		{
			name:      "nested non-string values",
			value:     map[string]any{"cores": 96},
			converted: true,
			expected:  map[string]string{"cores": "96"},
		},
		{
			name:      "empty map",
			value:     map[string]any{},
			converted: true,
			expected:  map[string]string{},
		},
		{
			name:      "not a map",
			value:     "not-a-map",
			converted: false,
		},
		{
			name:      "nil",
			value:     nil,
			converted: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result, converted := AsStringKeyedMap(testCase.value)
			assert.Equal(t, testCase.converted, converted)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestVerifyAPIVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		versions          oranapi.APIVersions
		expectedVersion   string
		expectedURIPrefix string
		wantErr           bool
	}{
		{
			name: "valid cluster versions",
			versions: oranapi.APIVersions{
				ApiVersions: &[]oranapi.APIVersion{{Version: new(tsparams.ClusterAPIVersion)}},
				UriPrefix:   new(tsparams.ClusterAPIURIPrefix),
			},
			expectedVersion:   tsparams.ClusterAPIVersion,
			expectedURIPrefix: tsparams.ClusterAPIURIPrefix,
		},
		{
			name: "valid inventory versions",
			versions: oranapi.APIVersions{
				ApiVersions: &[]oranapi.APIVersion{{Version: new(tsparams.InventoryAPIVersion)}},
				UriPrefix:   new(tsparams.InventoryAPIURIPrefix),
			},
			expectedVersion:   tsparams.InventoryAPIVersion,
			expectedURIPrefix: tsparams.InventoryAPIURIPrefix,
		},
		{
			name: "wrong version",
			versions: oranapi.APIVersions{
				ApiVersions: &[]oranapi.APIVersion{{Version: new(tsparams.InventoryAPIVersion)}},
				UriPrefix:   new(tsparams.ClusterAPIURIPrefix),
			},
			expectedVersion:   tsparams.ClusterAPIVersion,
			expectedURIPrefix: tsparams.ClusterAPIURIPrefix,
			wantErr:           true,
		},
		{
			name: "wrong uriPrefix",
			versions: oranapi.APIVersions{
				ApiVersions: &[]oranapi.APIVersion{{Version: new(tsparams.ClusterAPIVersion)}},
				UriPrefix:   new("/wrong"),
			},
			expectedVersion:   tsparams.ClusterAPIVersion,
			expectedURIPrefix: tsparams.ClusterAPIURIPrefix,
			wantErr:           true,
		},
		{
			name:              "missing fields",
			versions:          oranapi.APIVersions{},
			expectedVersion:   tsparams.ClusterAPIVersion,
			expectedURIPrefix: tsparams.ClusterAPIURIPrefix,
			wantErr:           true,
		},
		{
			name: "empty apiVersions",
			versions: oranapi.APIVersions{
				ApiVersions: &[]oranapi.APIVersion{},
				UriPrefix:   new(tsparams.ClusterAPIURIPrefix),
			},
			expectedVersion:   tsparams.ClusterAPIVersion,
			expectedURIPrefix: tsparams.ClusterAPIURIPrefix,
			wantErr:           true,
		},
		{
			name: "nil version",
			versions: oranapi.APIVersions{
				ApiVersions: &[]oranapi.APIVersion{{}},
				UriPrefix:   new(tsparams.ClusterAPIURIPrefix),
			},
			expectedVersion:   tsparams.ClusterAPIVersion,
			expectedURIPrefix: tsparams.ClusterAPIURIPrefix,
			wantErr:           true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := VerifyAPIVersions(testCase.versions, testCase.expectedVersion, testCase.expectedURIPrefix)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRefersTo(t *testing.T) {
	t.Parallel()

	const (
		objectID   = "pool-1"
		objectName = "test-pool"
	)

	tests := []struct {
		name      string
		objectRef *string
		post      *map[string]any
		prior     *map[string]any
		want      bool
	}{
		{
			name:      "objectRef contains id",
			objectRef: new("/resourcePools/" + objectID),
			want:      true,
		},
		{
			name: "id in post state",
			post: new(map[string]any{"resourcePoolId": objectID, "name": objectName}),
			want: true,
		},
		{
			name:  "id in prior state",
			prior: new(map[string]any{"resourcePoolId": objectID}),
			want:  true,
		},
		{
			name: "name in post state",
			post: new(map[string]any{"name": objectName}),
			want: true,
		},
		{
			name:      "objectRef does not contain id",
			objectRef: new("/resourcePools/other"),
		},
		{
			name: "all nil",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := RefersTo(
				testCase.objectRef, testCase.post, testCase.prior, "resourcePoolId", objectID, "name", objectName)
			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestExtensionString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		extensions map[string]any
		key        string
		want       string
	}{
		{
			name: "nil extensions",
			key:  tsparams.ClusterModelExtension,
		},
		{
			name:       "present string",
			extensions: map[string]any{tsparams.ClusterModelExtension: tsparams.ClusterModelHubCluster},
			key:        tsparams.ClusterModelExtension,
			want:       tsparams.ClusterModelHubCluster,
		},
		{
			name:       "missing key",
			extensions: map[string]any{tsparams.ClusterModelExtension: tsparams.ClusterModelHubCluster},
			key:        "missing",
		},
		{
			name:       "nil value",
			extensions: map[string]any{tsparams.ClusterModelExtension: nil},
			key:        tsparams.ClusterModelExtension,
		},
		{
			name:       "non-string value",
			extensions: map[string]any{tsparams.ClusterModelExtension: 4},
			key:        tsparams.ClusterModelExtension,
			want:       "4",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.want, ExtensionString(testCase.extensions, testCase.key))
		})
	}
}

func TestAppendMismatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		errs      []error
		field     string
		want      any
		got       any
		wantCount int
	}{
		{
			name:  "matching values",
			field: "name",
			want:  "a",
			got:   "a",
		},
		{
			name:      "mismatch appends",
			field:     "name",
			want:      "a",
			got:       "b",
			wantCount: 1,
		},
		{
			name:      "preserves existing errors",
			errs:      []error{fmt.Errorf("existing")},
			field:     "name",
			want:      "a",
			got:       "b",
			wantCount: 2,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := AppendMismatch(testCase.errs, testCase.field, testCase.want, testCase.got)
			assert.Len(t, result, testCase.wantCount)
		})
	}
}

func TestAppendError(t *testing.T) {
	t.Parallel()

	existing := fmt.Errorf("existing")
	newErr := fmt.Errorf("new")

	tests := []struct {
		name      string
		errs      []error
		err       error
		wantCount int
	}{
		{
			name: "nil error",
		},
		{
			name:      "appends error",
			err:       newErr,
			wantCount: 1,
		},
		{
			name:      "preserves existing errors",
			errs:      []error{existing},
			err:       newErr,
			wantCount: 2,
		},
		{
			name:      "nil error keeps existing",
			errs:      []error{existing},
			wantCount: 1,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := AppendError(testCase.errs, testCase.err)
			assert.Len(t, result, testCase.wantCount)
		})
	}
}
