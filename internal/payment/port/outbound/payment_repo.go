package outbound

import (
	"github.com/Flook0147/netflix_like/internal/payment/domain"
	"github.com/google/uuid"
)

type PaymentRepoPort interface {
	CreatePayment(viewerId uuid.UUID, amount float64, method string) error
	GetAllPayment(viewerId uuid.UUID) ([]*domain.Payment, error)
	GetPaymentById(paymentId uuid.UUID) (*domain.Payment, error)
	GetSinglePayment(viewerId uuid.UUID, amount float64, method string) (*domain.Payment, error)
	UpdatePaymentStatus(paymentID uuid.UUID, status string) error
}
