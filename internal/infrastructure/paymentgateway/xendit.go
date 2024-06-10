package paymentgateway

import (
	"context"
	"os"
	"payment-service/internal/domain"
	"strings"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/ewallet"
	"github.com/xendit/xendit-go/invoice"
	"github.com/xendit/xendit-go/qrcode"
	"github.com/xendit/xendit-go/virtualaccount"
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
		ExternalID: payment.PaymentID,
		Amount:     payment.Amount,
		Currency:   payment.Currency,
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

func (xc *XenditClient) ChargeEWallet(ctx context.Context, payment *domain.Payment) (string, error) {
	xendit.Opt.SecretKey = xc.apiKey
	params := ewallet.CreateEWalletChargeParams{
		ReferenceID:    payment.PaymentID,
		Currency:       payment.Currency,
		Amount:         float64(payment.Amount * 100), // Xendit expects amount in cents
		CheckoutMethod: payment.EwalletCheckoutMethod,
	}

	charge, err := ewallet.CreateEWalletCharge(&params)
	if err != nil {
		return "", err
	}

	return charge.ID, nil
}

func (xc *XenditClient) CreateVirtualAccount(ctx context.Context, payment *domain.Payment) (string, error) {
	xendit.Opt.SecretKey = xc.apiKey

	// Remove the "XEN-" prefix from the payment method
	bankCode := strings.TrimPrefix(payment.PaymentMethod, "XEN-")

	params := virtualaccount.CreateFixedVAParams{
		ExternalID:     payment.PaymentID,
		BankCode:       bankCode,                      // Bank code, e.g., "BCA", "BNI", etc.
		Name:           payment.UserID,                // Assuming UserID is the name here
		ExpectedAmount: float64(payment.Amount * 100), // Xendit expects amount in cents
	}

	va, err := virtualaccount.CreateFixedVA(&params)
	if err != nil {
		return "", err
	}

	return va.ID, nil
}

func (xc *XenditClient) CreateQRCode(ctx context.Context, payment *domain.Payment) (string, error) {
	xendit.Opt.SecretKey = xc.apiKey

	params := qrcode.CreateQRCodeParams{
		ExternalID:  payment.PaymentID,
		Amount:      float64(payment.Amount * 100), // Xendit expects amount in cents
		Type:        xendit.QRCodeType(payment.QrType),
		CallbackURL: payment.QrCallbackURL,
	}

	qrCode, err := qrcode.CreateQRCode(&params)
	if err != nil {
		return "", err
	}

	return qrCode.QRString, nil
}
