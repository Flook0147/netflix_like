package db

import (
	"github.com/Flook0147/netflix_like/internal/user/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserHandler {
	return &UserHandler{
		DB: db,
	}
}

type Viewer struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid"`
}

func (h *UserHandler) CreateUser(username, password, name, email string) error {
	tx := h.DB.Begin()

	user := &domain.User{
		UserId:     uuid.New(), // 🔥 ต้องมี ID
		Username:   username,
		Password:   password,
		Name:       name,
		Email:      email,
		ProfileURL: "",
	}

	// 1. create user
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. create viewer
	viewer := &domain.Viewer{
		Viewer_id: uuid.New(),
		User_id:   user.UserId, // 🔥 link กับ user
	}

	if err := tx.Create(viewer).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (h *UserHandler) GetUser(username string) (*domain.User, error) {
	var user domain.User
	result := h.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
