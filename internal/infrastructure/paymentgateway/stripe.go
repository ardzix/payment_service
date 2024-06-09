// internal/infrastructure/paymentgateway/stripe.go
package paymentgateway

import (
	"context"
	"os"
	"payment-service/internal/domain"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/refund"
)

type StripeClient struct {
	apiKey string
}

func NewStripeClient() *StripeClient {
	return &StripeClient{
		apiKey: os.Getenv("STRIPE_API_KEY"),
	}
}

func (sc *StripeClient) ProcessPayment(ctx context.Context, payment *domain.Payment) (string, error) {
	stripe.Key = sc.apiKey

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(payment.Amount * 100)), // Stripe accepts amounts in cents
		Currency: stripe.String(payment.Currency),
		PaymentMethodTypes: stripe.StringSlice([]string{
			payment.PaymentMethod,
		}),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", err
	}

	return pi.ID, nil
}

func (sc *StripeClient) RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error) {
	stripe.Key = sc.apiKey

	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentID),
		Amount:        stripe.Int64(int64(amount * 100)),
	}

	refund, err := refund.New(params)
	if err != nil {
		return "", err
	}

	return refund.ID, nil
}
