package paymentgateway

import (
	"context"
	"log"
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

	log.Printf("Sending request to Xendit to create an invoice: %+v\n", data)
	createdInvoice, err := invoice.Create(&data)
	if err != nil {
		log.Printf("Error creating invoice with Xendit: %v\n", err)
		return "", err
	}

	log.Printf("Invoice created successfully with ID: %s\n", createdInvoice.ID)
	return createdInvoice.ID, nil
}

func (xc *XenditClient) RefundPayment(ctx context.Context, paymentID string, amount float64) (string, error) {
	// Implement refund logic with Xendit if available, as Xendit primarily supports invoice-based payments
	return "", nil
}

func (xc *XenditClient) ChargeEWallet(ctx context.Context, payment *domain.Payment) (string, error) {
	xendit.Opt.SecretKey = xc.apiKey
	channelProperties := map[string]string{
		"mobile_number":        "+6285811144421",
		"success_redirect_url": "https://arnatech.id",
	}
	params := ewallet.CreateEWalletChargeParams{
		ReferenceID:       payment.PaymentID,
		Currency:          payment.Currency,
		Amount:            payment.Amount,
		CheckoutMethod:    payment.EwalletCheckoutMethod,
		ChannelCode:       "ID_" + payment.PaymentMethod,
		ChannelProperties: channelProperties,
	}

	log.Printf("Sending request to Xendit to charge e-wallet: %+v\n", params)
	charge, err := ewallet.CreateEWalletCharge(&params)
	if err != nil {
		log.Printf("Error charging e-wallet with Xendit: %v\n", err)
		return "", err
	}

	log.Printf("E-wallet charged successfully with ID: %s\n", charge.ID)
	return charge.ID, nil
}

func (xc *XenditClient) CreateVirtualAccount(ctx context.Context, payment *domain.Payment) (string, error) {
	xendit.Opt.SecretKey = xc.apiKey

	// Remove the "XEN-" prefix from the payment method
	bankCode := strings.TrimPrefix(payment.PaymentMethod, "XEN-")
	trueValue := true
	params := virtualaccount.CreateFixedVAParams{
		ExternalID:     payment.PaymentID,
		BankCode:       bankCode,       // Bank code, e.g., "BCA", "BNI", etc.
		Name:           payment.UserID, // Assuming UserID is the name here
		ExpectedAmount: payment.Amount,
		IsClosed:       &trueValue,
	}

	log.Printf("Sending request to Xendit to create virtual account: %+v\n", params)
	va, err := virtualaccount.CreateFixedVA(&params)
	if err != nil {
		log.Printf("Error creating virtual account with Xendit: %v\n", err)
		return "", err
	}

	log.Printf("Virtual account created successfully with ID: %s\n", va.ID)
	return va.ID, nil
}

func (xc *XenditClient) CreateQRCode(ctx context.Context, payment *domain.Payment) (string, error) {
	xendit.Opt.SecretKey = xc.apiKey

	params := qrcode.CreateQRCodeParams{
		ExternalID:  payment.PaymentID,
		Amount:      payment.Amount,
		Type:        xendit.QRCodeType(payment.QrType),
		CallbackURL: payment.QrCallbackURL,
	}

	log.Printf("Sending request to Xendit to create QR code: %+v\n", params)
	qrCode, err := qrcode.CreateQRCode(&params)
	if err != nil {
		log.Printf("Error creating QR code with Xendit: %v\n", err)
		return "", err
	}

	payment.QrString = qrCode.QRString

	log.Printf("QR code created successfully with details: %+v\n", qrCode)
	return qrCode.ID, nil
}
