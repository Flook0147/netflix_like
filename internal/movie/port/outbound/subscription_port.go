package outbound

import "github.com/google/uuid"

type SubscriptionPort interface {
	GetSubscriptionStatus(viewerID uuid.UUID) (string, error)
}
