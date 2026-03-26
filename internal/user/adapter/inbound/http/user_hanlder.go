package http

import (
	"strings"

	"github.com/Flook0147/netflix_like/internal/auth/port/inbound"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	authPort inbound.AuthPort
}

func NewUserHandler(authPort inbound.AuthPort) *UserHandler {
	return &UserHandler{
		authPort: authPort,
	}
}

type ValidateTokenRequest struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *UserHandler) GetUserFromToken(c fiber.Ctx) error {
	token := c.Get("Authorization")

	if token == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "missing token",
		})
	}

	token = strings.TrimPrefix(token, "Bearer ")

	user, err := h.authPort.GetUserFromToken(token) // ✅ FIX HERE
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"user": user,
	})
}
