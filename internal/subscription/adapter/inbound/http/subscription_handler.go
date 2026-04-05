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

// Subscribe godoc
// @Summary Subscribe to a plan
// @Description Subscribe current user to a plan
// @Tags subscription
// @Accept json
// @Produce json
// @Param body body object{plan_id=string} true "Plan ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /subscription [post]
// @Security BearerAuth
func (h *SubscriptionHandler) Subscribe(c fiber.Ctx) error {
	type reqBody struct {
		PlanID   string `json:"plan_id"`
		ViewerID string `json:"viewer_id"`
	}

	var body reqBody
	if err := c.Bind().Body(&body); err != nil {
		return err
	}

	planID, err := uuid.Parse(body.PlanID)
	if err != nil {
		return err
	}

	viewerID, err := uuid.Parse(body.ViewerID)
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

// GetMySubscription godoc
// @Summary Get my subscription
// @Description Get current user's subscription status
// @Tags subscription
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /subscription [get]
// @Security BearerAuth
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
