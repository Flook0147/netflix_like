package db

import (
	"time"

	"github.com/Flook0147/netflix_like/internal/payment/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

const (
	STATUSPENING = "PENDING"
)

func (p *PaymentRepository) CreatePayment(viewerId uuid.UUID, amount float64, method string) error {
	payment := domain.Payment{
		ViewerID:  viewerId,
		Amount:    amount,
		Status:    STATUSPENING,
		Method:    method,
		SlipURL:   "",
		CreatedAt: time.Now(),
	}
	return p.db.Create(payment).Error
}

func (p *PaymentRepository) GetAllPayment(viewerId uuid.UUID) ([]*domain.Payment, error) {
	var payment []*domain.Payment
	err := p.db.Where("viewer_id = ?", viewerId).Find(&payment).Error
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (p *PaymentRepository) GetSinglePayment(viewerId uuid.UUID, amount float64, method string) (*domain.Payment, error) {
	var payment *domain.Payment
	err := p.db.Where("viewer_id = ? AND amount = ? AND method = ?", viewerId, amount, method).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (p *PaymentRepository) GetPaymentById(paymentId uuid.UUID) (*domain.Payment, error) {
	var payment *domain.Payment
	err := p.db.Where("payment_id = ?", paymentId).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (p *PaymentRepository) UpdatePaymentStatus(paymentId uuid.UUID, status string) error {
	var payment *domain.Payment
	return p.db.Model(domain.Payment{}).
		Where("payment_id = ?", paymentId).First(&payment).
		Update("status", status).Error
}
