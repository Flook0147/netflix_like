package http

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, handler *UserHandler) {

	user := router.Group("/user")

	user.Get("/me", handler.GetUserFromToken)
}
