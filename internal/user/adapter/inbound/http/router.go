package http

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App, handler *UserHandler) {

	user := app.Group("/user")

	user.Get("/me", handler.GetUserFromToken)
}
