package db

import (
	"github.com/Flook0147/netflix_like/internal/user/domain"
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

func (h *UserHandler) CreateUser(username, password, name, email string) error {
	// Implement user creation logic here, e.g., save to database
	result := h.DB.Create(&domain.User{
		Username:   username,
		Password:   password,
		Name:       name,
		Email:      email,
		ProfileURL: "",
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (h *UserHandler) GetUser(username string) (*domain.User, error) {
	var user domain.User
	result := h.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
