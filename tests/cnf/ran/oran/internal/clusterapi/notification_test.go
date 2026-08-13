package clusterapi

import (
	"testing"

	"github.com/google/uuid"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"github.com/stretchr/testify/assert"
)

func TestMatchModifyNodeCluster(t *testing.T) {
	t.Parallel()

	const (
		nodeClusterName = "spoke-cluster"
		labelValue      = "test-label-value"
	)

	nodeClusterID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	subscriptionID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	match := MatchModifyNodeCluster(
		subscriptionID,
		nodeClusterID,
		nodeClusterName,
		tsparams.TestNotificationLabel,
		labelValue,
	)

	postObjectState := map[string]any{
		"nodeClusterId": nodeClusterID.String(),
		"name":          nodeClusterName,
		"extensions": map[string]any{
			tsparams.TestNotificationLabel: labelValue,
		},
	}
	priorObjectState := map[string]any{"nodeClusterId": nodeClusterID.String()}

	notification := &oranapi.ClusterChangeNotification{
		NotificationEventType:  oranapi.ClusterChangeNotificationEventTypeModify,
		ConsumerSubscriptionId: &subscriptionID,
		PostObjectState:        &postObjectState,
		PriorObjectState:       &priorObjectState,
	}

	assert.True(t, match(notification))

	assert.False(t, match(nil))
	assert.False(t, match(&oranapi.ClusterChangeNotification{
		NotificationEventType:  oranapi.ClusterChangeNotificationEventTypeCreate,
		ConsumerSubscriptionId: &subscriptionID,
	}))
	assert.False(t, match(&oranapi.ClusterChangeNotification{
		NotificationEventType:  oranapi.ClusterChangeNotificationEventTypeModify,
		ConsumerSubscriptionId: &subscriptionID,
		PostObjectState:        &postObjectState,
	}))

	wrongLabelPost := map[string]any{
		"nodeClusterId": nodeClusterID.String(),
		"name":          nodeClusterName,
		"extensions": map[string]any{
			tsparams.TestNotificationLabel: "other-value",
		},
	}
	assert.False(t, match(&oranapi.ClusterChangeNotification{
		NotificationEventType:  oranapi.ClusterChangeNotificationEventTypeModify,
		ConsumerSubscriptionId: &subscriptionID,
		PostObjectState:        &wrongLabelPost,
		PriorObjectState:       &priorObjectState,
	}))

	missingExtensionsPost := map[string]any{
		"nodeClusterId": nodeClusterID.String(),
		"name":          nodeClusterName,
	}
	assert.False(t, match(&oranapi.ClusterChangeNotification{
		NotificationEventType:  oranapi.ClusterChangeNotificationEventTypeModify,
		ConsumerSubscriptionId: &subscriptionID,
		PostObjectState:        &missingExtensionsPost,
		PriorObjectState:       &priorObjectState,
	}))
}
