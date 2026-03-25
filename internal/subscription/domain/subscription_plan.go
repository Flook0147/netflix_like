package domain

import (
	"github.com/google/uuid"
)

type SubscriptionPlan struct {
	ID        uuid.UUID
	Name      string
	Price     float64
	ValidDays int
}
