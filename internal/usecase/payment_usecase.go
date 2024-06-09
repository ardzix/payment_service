// internal/usecase/payment_usecase.go
package usecase

import (
	"context"
	"errors"
	"payment-service/internal/domain"
)

type paymentUseCase struct {
	paymentRepo  domain.PaymentRepository
	stripeClient domain.PaymentGateway
	xenditClient domain.PaymentGateway
}

func NewPaymentUseCase(repo domain.PaymentRepository, stripeClient, xenditClient domain.PaymentGateway) PaymentUseCase {
	return &paymentUseCase{
		paymentRepo:  repo,
		stripeClient: stripeClient,
		xenditClient: xenditClient,
	}
}

func (uc *paymentUseCase) ProcessPayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	var paymentID string
	var err error

	switch payment.PaymentMethod {
	case "stripe":
		paymentID, err = uc.stripeClient.ProcessPayment(ctx, payment)
	case "xendit":
		paymentID, err = uc.xenditClient.ProcessPayment(ctx, payment)
	default:
		return nil, errors.New("unsupported payment method")
	}

	if err != nil {
		return nil, err
	}

	payment.PaymentID = paymentID
	payment.Status = "processed"

	err = uc.paymentRepo.Save(payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (uc *paymentUseCase) RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error) {
	// Implement refund logic with appropriate gateway
	return "", nil
}

func (uc *paymentUseCase) GetPaymentStatus(ctx context.Context, paymentID string) (*domain.Payment, error) {
	payment, err := uc.paymentRepo.FindByID(paymentID)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (uc *paymentUseCase) ListPayments(ctx context.Context, userID string, page, pageSize int) ([]domain.Payment, int, error) {
	payments, totalCount, err := uc.paymentRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return payments, totalCount, nil
}
