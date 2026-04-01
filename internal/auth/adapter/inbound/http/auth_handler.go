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

// Register godoc
// @Summary Register user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Register request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
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

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access + refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
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

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RefreshRequest true "Refresh token"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {

	var req RefreshRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "missing refresh token",
		})
	}

	newToken, newRefreshToken, err := h.auth.RefreshToken(req.RefreshToken)
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
