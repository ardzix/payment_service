// internal/usecase/usecase.go
package usecase

import (
	"context"
	"payment-service/internal/domain"
)

type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error)
	GetPaymentStatus(ctx context.Context, paymentID string) (*domain.Payment, error)
	ListPayments(ctx context.Context, userID string, page, pageSize int) ([]domain.Payment, int, error)
}
