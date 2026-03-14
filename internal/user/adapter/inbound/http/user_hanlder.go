package http

import (
	"github.com/Flook0147/netflix_like/internal/user/port/inbound"
)

type UserHandler struct {
	userPort inbound.UserPort
}

func NewUserHandler(userPort inbound.UserPort) *UserHandler {
	return &UserHandler{
		userPort: userPort,
	}
}

type ValidateTokenRequest struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
