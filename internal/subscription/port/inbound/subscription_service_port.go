package inbound

import (
	"github.com/Flook0147/netflix_like/internal/subscription/domain"
	"github.com/google/uuid"
)

type SubscriptionServicePort interface {
	Subscribe(viewerID, planID uuid.UUID) error
	Unsubscribe(viewerID, planID uuid.UUID) error
	GetSubscriptionStatus(viewerID, planID uuid.UUID) (string, error)
	ActivateSubscriptions(viewerID, planID uuid.UUID) error
	GetSubscriptionDetails(viewerID uuid.UUID) ([]*domain.Subscription, error)
}
