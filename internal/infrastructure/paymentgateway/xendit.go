// internal/infrastructure/paymentgateway/xendit.go
package paymentgateway

import (
	"context"
	"os"
	"payment-service/internal/domain"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)

type XenditClient struct {
	apiKey string
}

func NewXenditClient() *XenditClient {
	return &XenditClient{
		apiKey: os.Getenv("XENDIT_API_KEY"),
	}
}

func (xc *XenditClient) ProcessPayment(ctx context.Context, payment *domain.Payment) (string, error) {
	xendit.Opt.SecretKey = xc.apiKey

	data := invoice.CreateParams{
		ExternalID:  payment.PaymentID,
		Amount:      payment.Amount,
		PayerEmail:  payment.UserID, // Assuming UserID is an email
		Description: "Payment for Order",
		Currency:    payment.Currency,
		PaymentMethods: []string{
			payment.PaymentMethod,
		},
	}

	createdInvoice, err := invoice.Create(&data)
	if err != nil {
		return "", err
	}

	return createdInvoice.ID, nil
}

func (xc *XenditClient) RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error) {
	// Implement refund logic with Xendit if available, as Xendit primarily supports invoice-based payments
	return "", nil
}
