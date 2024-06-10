// internal/domain/repository.go
package domain

import "context"

type PaymentRepository interface {
	Save(ctx context.Context, payment *Payment) error
	FindByID(ctx context.Context, paymentID string) (*Payment, error)
	FindByUserID(ctx context.Context, userID string, page, pageSize int) ([]Payment, int, error)
	UpdateStatus(ctx context.Context, paymentID, status string) error
}

type PaymentGateway interface {
	ProcessPayment(ctx context.Context, payment *Payment) (string, error)
	RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error)
	ChargeEWallet(ctx context.Context, payment *Payment) (string, error)
	CreateVirtualAccount(ctx context.Context, payment *Payment) (string, error)
	CreateQRCode(ctx context.Context, payment *Payment) (string, error)
}
