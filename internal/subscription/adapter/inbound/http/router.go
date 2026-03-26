package http

import "github.com/gofiber/fiber/v3"

func RegisterSubscriptionRoutes(app fiber.Router, handler *SubscriptionHandler) {
	sub := app.Group("/subscriptions")

	sub.Post("/subscribe", handler.Subscribe)
	sub.Get("/me", handler.GetMySubscription)
}
