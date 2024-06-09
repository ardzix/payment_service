// internal/domain/entity.go
package domain

type Payment struct {
    PaymentID string
    UserID    string
    Amount    float64
    Currency  string
    Status    string
    CreatedAt string
    PaymentMethod string
}

