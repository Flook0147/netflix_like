package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID        uuid.UUID `json:"id"`
	ViewerID  uuid.UUID `json:"viewer_id"`
	PlanID    uuid.UUID `json:"plan_id"`
	PaymentID uuid.UUID `json:"payment_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
}
