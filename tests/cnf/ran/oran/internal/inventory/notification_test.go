package inventory

import (
	"testing"

	"github.com/google/uuid"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"github.com/stretchr/testify/assert"
)

func TestMatchResourcePoolChange(t *testing.T) {
	t.Parallel()

	const poolID = "abc-123-pool-id"

	subscriptionID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	match := MatchResourcePoolChange(
		oranapi.InventoryChangeNotificationEventTypeCreate,
		subscriptionID,
		poolID,
		tsparams.TestInventoryResourcePool,
	)

	postObjectState := map[string]any{
		"resourcePoolId": poolID,
		"name":           tsparams.TestInventoryResourcePool,
	}

	notification := &oranapi.InventoryChangeNotification{
		NotificationEventType:  oranapi.InventoryChangeNotificationEventTypeCreate,
		ConsumerSubscriptionId: &subscriptionID,
		PostObjectState:        &postObjectState,
	}

	assert.True(t, match(notification))
	assert.True(t, match(&oranapi.InventoryChangeNotification{
		NotificationEventType:  oranapi.InventoryChangeNotificationEventTypeCreate,
		ConsumerSubscriptionId: &subscriptionID,
		ObjectRef:              new("/resourcePools/" + poolID),
	}))

	assert.False(t, match(nil))
	assert.False(t, match(&oranapi.InventoryChangeNotification{
		NotificationEventType:  oranapi.InventoryChangeNotificationEventTypeDelete,
		ConsumerSubscriptionId: &subscriptionID,
		PostObjectState:        &postObjectState,
	}))
}
