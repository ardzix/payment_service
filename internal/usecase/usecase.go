// internal/usecase/usecase.go
package usecase

import (
	"context"
	"errors"
	"payment-service/internal/domain"
	"payment-service/internal/infrastructure/paymentgateway"
)

type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error)
	GetPaymentStatus(ctx context.Context, paymentID string) (*domain.Payment, error)
	ListPayments(ctx context.Context, userID string, page, pageSize int) ([]domain.Payment, int, error)
}

type paymentUseCase struct {
	stripeClient *paymentgateway.StripeClient
	xenditClient *paymentgateway.XenditClient
	paymentRepo  domain.PaymentRepository
}

func NewPaymentUseCase(stripeClient *paymentgateway.StripeClient, xenditClient *paymentgateway.XenditClient, paymentRepo domain.PaymentRepository) PaymentUseCase {
	return &paymentUseCase{
		stripeClient: stripeClient,
		xenditClient: xenditClient,
		paymentRepo:  paymentRepo,
	}
}

func (uc *paymentUseCase) ProcessPayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	var paymentID string
	var err error

	switch payment.PaymentMethod {
	case "XEN-OVO", "XEN-DANA", "XEN-LINKAJA":
		paymentID, err = uc.xenditClient.ChargeEWallet(ctx, payment)
	case "XEN-BCA", "XEN-BNI", "XEN-BRI":
		paymentID, err = uc.xenditClient.CreateVirtualAccount(ctx, payment)
	case "XEN-QR":
		paymentID, err = uc.xenditClient.CreateQRCode(ctx, payment)
	case "XEN-DEFAULT":
		paymentID, err = uc.xenditClient.ProcessPayment(ctx, payment) // Xendit default
	case "STP-DEFAULT":
		paymentID, err = uc.stripeClient.ProcessPayment(ctx, payment) // Stripe default
	default:
		return nil, errors.New("unsupported payment method")
	}

	if err != nil {
		return nil, err
	}

	payment.PaymentID = paymentID
	payment.Status = "pending"
	err = uc.paymentRepo.Save(ctx, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (uc *paymentUseCase) RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error) {
	// Retrieve payment to determine which client to use for refund
	payment, err := uc.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return "", err
	}

	var refundID string
	switch payment.PaymentMethod {
	case "XEN-OVO", "XEN-DANA", "XEN-LINKAJA", "XEN-BCA", "XEN-BNI", "XEN-BRI", "XEN-QR", "XEN-DEFAULT":
		refundID, err = uc.xenditClient.RefundPayment(ctx, paymentID, amount)
	case "STP-DEFAULT":
		refundID, err = uc.stripeClient.RefundPayment(ctx, paymentID, amount)
	default:
		return "", errors.New("unsupported payment method")
	}

	if err != nil {
		return "", err
	}

	err = uc.paymentRepo.UpdateStatus(ctx, paymentID, "refunded")
	if err != nil {
		return "", err
	}

	return refundID, nil
}

func (uc *paymentUseCase) GetPaymentStatus(ctx context.Context, paymentID string) (*domain.Payment, error) {
	payment, err := uc.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (uc *paymentUseCase) ListPayments(ctx context.Context, userID string, page, pageSize int) ([]domain.Payment, int, error) {
	payments, total, err := uc.paymentRepo.FindByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return payments, total, nil
}
