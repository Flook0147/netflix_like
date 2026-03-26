package http

import (
	"github.com/Flook0147/netflix_like/internal/subscription/port/inbound"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	service inbound.SubscriptionPort
}

func NewSubscriptionHandler(service inbound.SubscriptionPort) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) Subscribe(c fiber.Ctx) error {
	type reqBody struct {
		PlanID string `json:"plan_id"`
	}

	var body reqBody
	if err := c.Bind().Body(&body); err != nil {
		return err
	}

	viewerID := c.Locals("userID").(uuid.UUID)

	planID, err := uuid.Parse(body.PlanID)
	if err != nil {
		return err
	}

	err = h.service.Subscribe(viewerID, planID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "subscribed successfully",
	})
}

func (h *SubscriptionHandler) GetMySubscription(c fiber.Ctx) error {
	viewerID := c.Locals("userID").(uuid.UUID)

	status, err := h.service.GetSubscriptionStatus(viewerID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": status,
	})
}
