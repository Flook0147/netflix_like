package outbound

import (
	"time"

	"github.com/Flook0147/netflix_like/internal/subscription/domain"
	"github.com/google/uuid"
)

type SubscriptionRepoPort interface {
	CreateSubscription(viwerId, planId uuid.UUID) error
	GetAllSubscription(viewerId uuid.UUID) ([]*domain.Subscription, error)
	GetSubscription(viewerId uuid.UUID) (*domain.Subscription, error)
	UpdateSubscription(viewerId, planId uuid.UUID, status string, expired time.Time) error
}

type SubscriptionPlanPort interface {
	GetSubscriptionPlanDetail(planId uuid.UUID) (*domain.SubscriptionPlan, error)
}
