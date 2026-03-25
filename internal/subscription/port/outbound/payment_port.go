package outbound

import (
	"github.com/google/uuid"
)

type PaymentPort interface {
	VerifyPayment(paymentID uuid.UUID) (bool, error)
}
