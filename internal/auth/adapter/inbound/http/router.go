package http

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app fiber.Router, handler *AuthHandler) {

	auth := app.Group("/auth")

	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Post("/refresh", handler.RefreshToken)
}
