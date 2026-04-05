package domain

import "github.com/google/uuid"

type User struct {
	UserId     uuid.UUID `gorm:"type:uuid;default:uuid_generate();primaryKey"`
	Username   string    `gorm:"unique;not null"`
	Password   string    `gorm:"not null"`
	Name       string    `gorm:"not null"`
	Email      string    `gorm:"unique;not null"`
	ProfileURL string
	Role       string `gorm:"not null;default:'user'"`
}

type Viewer struct {
	User_id   uuid.UUID `gorm:"type:uuid;"`
	Viewer_id uuid.UUID `gorm:"type:uuid;default:uuid_generate();primaryKey"`
}

type Admin struct {
	User_id  uuid.UUID `gorm:"type:uuid;"`
	Admin_id uuid.UUID `gorm:"type:uuid;default:uuid_generate();primaryKey"`
}
