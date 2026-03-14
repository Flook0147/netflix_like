package domain

import "github.com/google/uuid"

type User struct {
	UserId     uuid.UUID `gorm:"type:uuid;default:uuid_generate();primaryKey"`
	Username   string    `gorm:"unique;not null"`
	Password   string    `gorm:"not null"`
	Name       string    `gorm:"not null"`
	Email      string    `gorm:"unique;not null"`
	ProfileURL string
}

type Viewer struct {
	User
}

type Administrator struct {
	User
}
