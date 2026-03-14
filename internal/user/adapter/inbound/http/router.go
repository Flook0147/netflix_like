package http

import (
	"github.com/gofiber/fiber/v3"
)

func NewRouter(userHandler *UserHandler) *fiber.App {
	app := fiber.New()

	// Define user-related routes here, e.g.:
	app.Get("/user/me", userHandler.userPort.GetUserFromToken)
	return app
}
