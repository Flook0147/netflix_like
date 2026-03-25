package db

import (
	"github.com/Flook0147/netflix_like/internal/subscription/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionPlanRepo struct {
	db *gorm.DB
}

func NewSubscriptionPlanRepo(db *gorm.DB) *SubscriptionPlanRepo {
	return &SubscriptionPlanRepo{
		db: db,
	}
}

func (s *SubscriptionPlanRepo) GetSubscriptionPlanDetail(planId uuid.UUID) (*domain.SubscriptionPlan, error) {
	var plan domain.SubscriptionPlan
	err := s.db.Where("plan_id = ?", planId).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}
