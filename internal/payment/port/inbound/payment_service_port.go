package inbound

import (
	"github.com/Flook0147/netflix_like/internal/payment/domain"
	"github.com/google/uuid"
)

type PaymentServicePort interface {
	VerifyPayment(paymentID uuid.UUID) (bool, error)
	CreatePayment(viewerID uuid.UUID, amount float64, method string) error
	GetPaymentByID(paymentID uuid.UUID) (*domain.Payment, error)
	GetPaymentsByViewerId(viewerId uuid.UUID) ([]*domain.Payment, error)
}
