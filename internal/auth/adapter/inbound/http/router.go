package http

import (
	"github.com/gofiber/fiber/v3"
)

func NewRouter(authHandler *AuthHandler) *fiber.App {
	app := fiber.New()

	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)
	app.Post("/refresh", authHandler.RefreshToken)

	return app
}
