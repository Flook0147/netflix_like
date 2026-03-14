package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	TokenHash string    `gorm:"not null"`
	ExpiredAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	CreatedAt time.Time
}
