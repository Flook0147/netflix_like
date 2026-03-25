package app

import (
	"fmt"

	"github.com/Flook0147/netflix_like/internal/payment/domain"
	"github.com/Flook0147/netflix_like/internal/payment/port/outbound"
	"github.com/google/uuid"
)

const (
	StatusPending  = "PENDING"
	StatusSuccess  = "SUCCESS"
	StatusFailed   = "FAILED"
	StatusRefunded = "REFUNDED"
)

type PaymentService struct {
	paymentRepo outbound.PaymentRepoPort
}

func NewPaymentService(paymentRepo outbound.PaymentRepoPort) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
	}
}

func (s *PaymentService) VerifyPayment(paymentID uuid.UUID) (bool, error) {
	// 1. get payment
	payment, err := s.paymentRepo.GetPaymentById(paymentID)
	if err != nil {
		return false, err
	}

	// 2. check state
	if payment.Status != StatusPending {
		return false, fmt.Errorf("payment is not pending")
	}

	fmt.Println("🔍 Verifying payment...")

	// 3. simulate verify (later: OCR / gateway)
	success := true

	if !success {
		err = s.paymentRepo.UpdatePaymentStatus(paymentID, StatusFailed)
		return false, err
	}

	// 4. update status → SUCCESS
	err = s.paymentRepo.UpdatePaymentStatus(paymentID, StatusSuccess)
	if err != nil {
		return false, err
	}

	fmt.Println("✅ Payment verified & event published")

	return true, nil
}

// simple crud methods for payment records

func (s *PaymentService) CreatePayment(viewerID uuid.UUID, amount float64, method string) error {
	return s.paymentRepo.CreatePayment(viewerID, amount, method)
}

func (s *PaymentService) GetPaymentByID(paymentID uuid.UUID) (*domain.Payment, error) {
	return s.paymentRepo.GetPaymentById(paymentID)
}

func (s *PaymentService) GetPaymentsByViewerId(viewerId uuid.UUID) ([]*domain.Payment, error) {
	return s.paymentRepo.GetAllPayment(viewerId)
}

func (s *PaymentService) UpdatePaymentStatus(paymentID uuid.UUID, status string) error {
	_, err := s.paymentRepo.GetPaymentById(paymentID)
	if err != nil {
		return err
	}

	return s.paymentRepo.UpdatePaymentStatus(paymentID, status)
}
