package http

import "github.com/gofiber/fiber/v3"

func RegisterSubscriptionRoutes(router fiber.Router, handler *SubscriptionHandler) {
	sub := router.Group("/subscriptions")

	sub.Post("/subscribe", handler.Subscribe)
	sub.Get("/me", handler.GetMySubscription)
}
