// internal/domain/repository.go
package domain

type PaymentRepository interface {
    Save(payment *Payment) error
    FindByID(paymentID string) (*Payment, error)
    FindByUserID(userID string, page, pageSize int) ([]Payment, int, error)
    UpdateStatus(paymentID, status string) error
}

type PaymentGateway interface {
    ProcessPayment(ctx context.Context, payment *Payment) (string, error)
    RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error)
}

