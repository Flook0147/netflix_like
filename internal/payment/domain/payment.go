package domain

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID        uuid.UUID `json:"id"`
	ViewerID  uuid.UUID `json:"user_id"`
	PlanID    uuid.UUID
	Status    string    `json:"status"`
	SlipURL   string    `json:"slip_url"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	CreatedAt time.Time `json:"created_at"`
}
