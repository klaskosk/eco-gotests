package clusterapi

import (
	"encoding/json"
	"testing"

	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsStringKeyedMap(t *testing.T) {
	t.Parallel()

	fromAny, ok := asStringKeyedMap(map[string]any{"cores": "96", "architecture": "x86_64"})
	assert.True(t, ok)
	assert.Equal(t, "96", fromAny["cores"])
	assert.Equal(t, "x86_64", fromAny["architecture"])

	fromString, ok := asStringKeyedMap(map[string]string{"GiB": "256"})
	assert.True(t, ok)
	assert.Equal(t, "256", fromString["GiB"])

	_, ok = asStringKeyedMap("not-a-map")
	assert.False(t, ok)
}

func TestVerifyAPIVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "valid versions",
			json: `{"apiVersions":[{"version":"1.0.0"}],"uriPrefix":"/o2ims-infrastructureCluster/v1"}`,
		},
		{
			name:    "wrong version",
			json:    `{"apiVersions":[{"version":"2.0.0"}],"uriPrefix":"/o2ims-infrastructureCluster/v1"}`,
			wantErr: true,
		},
		{
			name:    "wrong uriPrefix",
			json:    `{"apiVersions":[{"version":"1.0.0"}],"uriPrefix":"/wrong"}`,
			wantErr: true,
		},
		{
			name:    "missing fields",
			json:    `{}`,
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var versions oranapi.APIVersions
			require.NoError(t, json.Unmarshal([]byte(testCase.json), &versions))

			err := VerifyAPIVersions(versions)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNotificationHasExtensionLabel(t *testing.T) {
	t.Parallel()

	notification := &oranapi.ClusterChangeNotification{
		PostObjectState: &map[string]any{
			"extensions": map[string]any{
				"oran-test-notification": "abc",
				"capacity":               map[string]any{"cpu": "4"},
			},
		},
	}

	assert.True(t, NotificationHasExtensionLabel(notification, "oran-test-notification", "abc"))
	assert.False(t, NotificationHasExtensionLabel(notification, "oran-test-notification", "other"))
	assert.False(t, NotificationHasExtensionLabel(&oranapi.ClusterChangeNotification{}, "oran-test-notification", "abc"))
}

func TestExtensionString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "", ExtensionString(nil, "model"))

	extensions := map[string]any{"model": "hub-cluster"}
	assert.Equal(t, "hub-cluster", ExtensionString(&extensions, "model"))
	assert.Equal(t, "", ExtensionString(&extensions, "missing"))
}
