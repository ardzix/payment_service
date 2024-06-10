package usecase

import (
	"context"
	"errors"
	"os"
	"payment-service/internal/domain"
	"payment-service/internal/infrastructure/paymentgateway"
)

type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error)
	GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error)
	ListPayments(ctx context.Context, userID string, page, pageSize int) ([]domain.Payment, int, error)
}

type paymentUseCase struct {
	stripeClient        *paymentgateway.StripeClient
	xenditClient        *paymentgateway.XenditClient
	dokuClient          *paymentgateway.DokuClient
	paymentRepo         domain.PaymentRepository
	paymentConfigClient *paymentgateway.PaymentConfigClient
	defaultPG           string
}

func NewPaymentUseCase(stripeClient *paymentgateway.StripeClient, xenditClient *paymentgateway.XenditClient, dokuClient *paymentgateway.DokuClient, paymentRepo domain.PaymentRepository, paymentConfigClient *paymentgateway.PaymentConfigClient) PaymentUseCase {
	defaultPG := os.Getenv("DEFAULT_PG")
	return &paymentUseCase{
		stripeClient:        stripeClient,
		xenditClient:        xenditClient,
		dokuClient:          dokuClient,
		paymentRepo:         paymentRepo,
		paymentConfigClient: paymentConfigClient,
		defaultPG:           defaultPG,
	}
}

func (uc *paymentUseCase) ProcessPayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	var paymentID string
	var err error

	gateway, fallbackGateway, err := uc.paymentConfigClient.GetPaymentGatewayConfig(payment.PaymentMethod)
	if err != nil || gateway == "" {
		gateway = uc.defaultPG
	}

	payment.Gateway = gateway
	paymentMethod := payment.PaymentMethod

	switch gateway {
	case "XENDIT", "Xendit", "xendit":
		paymentID, err = uc.processWithXendit(ctx, paymentMethod, payment)
	case "DOKU", "Doku", "doku":
		paymentID, err = uc.processWithDoku(ctx, paymentMethod, payment)
	case "STRIPE", "Stripe", "stripe":
		paymentID, err = uc.stripeClient.ProcessPayment(ctx, payment)
	default:
		return nil, errors.New("unsupported payment gateway")
	}

	if err != nil && fallbackGateway != "" {
		payment.Gateway = fallbackGateway
		switch fallbackGateway {
		case "XENDIT", "Xendit", "xendit":
			paymentID, err = uc.processWithXendit(ctx, paymentMethod, payment)
		case "DOKU", "Doku", "doku":
			paymentID, err = uc.processWithDoku(ctx, paymentMethod, payment)
		case "STRIPE", "Stripe", "stripe":
			paymentID, err = uc.stripeClient.ProcessPayment(ctx, payment)
		}
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

func (uc *paymentUseCase) processWithXendit(ctx context.Context, paymentMethod string, payment *domain.Payment) (string, error) {
	switch paymentMethod {
	case "OVO", "DANA", "LINKAJA":
		return uc.xenditClient.ChargeEWallet(ctx, payment)
	case "BCA", "BNI", "BRI":
		return uc.xenditClient.CreateVirtualAccount(ctx, payment)
	case "QR":
		return uc.xenditClient.CreateQRCode(ctx, payment)
	case "DEFAULT":
		return uc.xenditClient.ProcessPayment(ctx, payment)
	default:
		return "", errors.New("unsupported payment method for Xendit")
	}
}

func (uc *paymentUseCase) processWithDoku(ctx context.Context, paymentMethod string, payment *domain.Payment) (string, error) {
	switch paymentMethod {
	case "OVO", "DANA", "LINKAJA":
		return uc.dokuClient.ChargeEWallet(ctx, payment)
	case "BCA", "BNI", "BRI":
		return uc.dokuClient.CreateVirtualAccount(ctx, payment)
	case "QR":
		return uc.dokuClient.CreateQRCode(ctx, payment)
	case "DEFAULT":
		return uc.dokuClient.ProcessPayment(ctx, payment)
	default:
		return "", errors.New("unsupported payment method for Doku")
	}
}

func (uc *paymentUseCase) RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error) {
	// Retrieve payment to determine which client to use for refund
	payment, err := uc.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return "", err
	}

	var refundID string
	switch payment.PaymentMethod {
	case "OVO", "DANA", "LINKAJA", "BCA", "BNI", "BRI", "QR", "DEFAULT":
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

func (uc *paymentUseCase) GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error) {
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
