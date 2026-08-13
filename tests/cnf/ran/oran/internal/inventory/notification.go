package inventory

import (
	"github.com/google/uuid"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/o2imstest"
)

// MatchResourcePoolChange returns a matcher for inventory notifications on the given ResourcePool.
//
// A notification matches when the event type and subscription ID agree, and the pool is referenced
// either in objectRef or via resourcePoolId/name in postObjectState or priorObjectState.
func MatchResourcePoolChange(
	event oranapi.InventoryChangeNotificationEventType,
	subscriptionID uuid.UUID,
	poolID, poolName string,
) func(*oranapi.InventoryChangeNotification) bool {
	return func(notification *oranapi.InventoryChangeNotification) bool {
		if notification == nil {
			return false
		}

		if notification.NotificationEventType != event {
			return false
		}

		if notification.ConsumerSubscriptionId == nil || *notification.ConsumerSubscriptionId != subscriptionID {
			return false
		}

		return o2imstest.RefersTo(
			notification.ObjectRef,
			notification.PostObjectState,
			notification.PriorObjectState,
			"resourcePoolId", poolID,
			"name", poolName,
		)
	}
}
