package paymentgateway

import (
	"context"
	"os"
	"payment-service/internal/domain"
)

type DokuClient struct {
	apiKey string
}

func NewDokuClient() *DokuClient {
	return &DokuClient{
		apiKey: os.Getenv("Doku_API_KEY"),
	}
}

func (xc *DokuClient) ProcessPayment(ctx context.Context, payment *domain.Payment) (string, error) {

	return payment.PaymentID, nil
}

func (xc *DokuClient) RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error) {
	// Implement refund logic with Doku if available, as Doku primarily supports invoice-based payments
	return "", nil
}

func (xc *DokuClient) ChargeEWallet(ctx context.Context, payment *domain.Payment) (string, error) {

	return payment.PaymentID, nil
}

func (xc *DokuClient) CreateVirtualAccount(ctx context.Context, payment *domain.Payment) (string, error) {

	return payment.PaymentID, nil
}

func (xc *DokuClient) CreateQRCode(ctx context.Context, payment *domain.Payment) (string, error) {

	return payment.PaymentID, nil
}
