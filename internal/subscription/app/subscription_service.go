package app

import (
	"fmt"
	"time"

	"github.com/Flook0147/netflix_like/internal/subscription/domain"
	"github.com/Flook0147/netflix_like/internal/subscription/port/outbound"
	"github.com/google/uuid"
)

type SubscriptionService struct {
	subRepo  outbound.SubscriptionRepoPort
	planRepo outbound.SubscriptionPlanPort
	payment  outbound.PaymentPort
}

const (
	StatusPending  = "PENDING"
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
	StatusFailed   = "FAILED"
)

func NewSubscriptionService(subRepo outbound.SubscriptionRepoPort, planRepo outbound.SubscriptionPlanPort, payment outbound.PaymentPort) *SubscriptionService {
	return &SubscriptionService{
		subRepo:  subRepo,
		planRepo: planRepo,
		payment:  payment,
	}
}

func (s *SubscriptionService) Subscribe(viewerID, planID uuid.UUID) error {
	sub, err := s.subRepo.GetSubscription(viewerID, planID)
	if err == nil && sub.Status != StatusInactive {
		return fmt.Errorf("already subscribed")
	}
	err = s.subRepo.CreateSubscription(viewerID, planID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) Unsubscribe(viewerID, planID uuid.UUID) error {
	sub, err := s.subRepo.GetSubscription(viewerID, planID)
	if err != nil {
		return err
	}

	if sub.Status != StatusActive {
		return fmt.Errorf("cannot unsubscribe")
	}

	return s.subRepo.UpdateSubscription(viewerID, planID, StatusInactive, time.Now())
}

func (s *SubscriptionService) GetSubscriptionStatus(viewerID, planID uuid.UUID) (string, error) {
	sub, err := s.subRepo.GetSubscription(viewerID, planID)
	if err != nil {
		return StatusFailed, err
	}
	return sub.Status, nil
}

func (s *SubscriptionService) ActivateSubscriptions(viewerID, planID uuid.UUID) error {
	sub, err := s.subRepo.GetSubscription(viewerID, planID)
	if err != nil {
		return err
	}
	if sub.Status != StatusPending {
		return fmt.Errorf("Invalid state transaction")
	}
	plan, err := s.planRepo.GetSubscriptionPlanDetail(planID)
	if err != nil {
		return err
	}

	valid, err := s.payment.VerifyPayment(sub.PaymentID)

	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("Payment is invalid")
	}

	return s.subRepo.UpdateSubscription(viewerID, planID, StatusActive, time.Now().AddDate(0, 0, plan.ValidDays))
}

func (s *SubscriptionService) GetSubscriptionDetails(viewerID uuid.UUID) ([]*domain.Subscription, error) {
	return s.subRepo.GetAllSubscription(viewerID)
}
