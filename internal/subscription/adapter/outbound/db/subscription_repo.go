package db

import (
	"time"

	"github.com/Flook0147/netflix_like/internal/subscription/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusPending  = "PENDING"
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
	StatusFailed   = "FAILED"
)

type SubscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepo(db *gorm.DB) *SubscriptionRepo {
	return &SubscriptionRepo{
		db: db,
	}
}

func (s *SubscriptionRepo) CreateSubscription(viwerId, planId uuid.UUID) error {
	sub := domain.Subscription{
		ViewerID: viwerId,
		PlanID:   planId,
		Status:   StatusPending,
	}
	return s.db.Create(&sub).Error
}

func (s *SubscriptionRepo) GetAllSubscription(viewerId uuid.UUID) ([]*domain.Subscription, error) {
	var sub []*domain.Subscription
	err := s.db.Where("viewer_id = ?", viewerId).Find(&sub).Error
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionRepo) GetSubscription(viewerId uuid.UUID) (*domain.Subscription, error) {
	var sub *domain.Subscription
	err := s.db.Where("viewer_id = ?", viewerId).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionRepo) UpdateSubscription(viewerId, planId uuid.UUID, status string, expired time.Time) error {
	return s.db.
		Model(&domain.Subscription{}).
		Where("viewer_id = ? AND plan_id = ?", viewerId, planId).
		Update("status", status).Error
}
