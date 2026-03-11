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

func (h *UserHandler) CreateUser(username, password, name string) error {
	// Implement user creation logic here, e.g., save to database
	err := h.DB.Create(&domain.User{
		Username:   username,
		Password:   password,
		Name:       name,
		ProfileURL: "",
	})
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (h *UserHandler) GetUser(username string) (*domain.User, error) {
	var user domain.User
	err := h.DB.Where("username = ?", username).First(&user)
	if err.Error != nil {
		return nil, err.Error
	}
	return &user, nil
}
