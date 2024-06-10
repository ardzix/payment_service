// internal/usecase/payment_gateway.go
package usecase

import (
	"context"
	"payment-service/internal/domain"
)

type PaymentGateway interface {
	ProcessPayment(ctx context.Context, payment *domain.Payment) (string, error)
	ChargeEWallet(ctx context.Context, payment *domain.Payment) (string, error)
	CreateVirtualAccount(ctx context.Context, payment *domain.Payment) (string, error)
	CreateQRCode(ctx context.Context, payment *domain.Payment) (string, error)
	RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error)
}
