package http

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, handler *AuthHandler) {

	auth := router.Group("/auth")

	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Post("/refresh", handler.RefreshToken)
}
