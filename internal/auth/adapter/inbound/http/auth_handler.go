package http

import (
	"github.com/Flook0147/netflix_like/internal/auth/port/inbound"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	auth inbound.AuthPort
}

func NewAuthHandler(auth inbound.AuthPort) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req RegisterRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	err := h.auth.Register(req.Username, req.Password, req.Name, req.Email)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "user created",
	})
}

func (h *AuthHandler) Login(c fiber.Ctx) error {

	var req LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	token, refreshToken, err := h.auth.Login(req.Username, req.Password)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}
	token = token[7:] // Remove "Bearer " prefix

	newToken, newRefreshToken, err := h.auth.RefreshToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(LoginResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
	})
}
