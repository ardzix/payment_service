package domain

import "time"

type Item struct {
	ItemName string
	Quantity int
	Price    float64
}

type Payment struct {
	PaymentID             string
	UserID                string
	Amount                float64
	Gateway               string
	Currency              string
	Status                string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	PaymentMethod         string
	PhoneNumber           string
	EwalletCheckoutMethod string
	QrType                string
	QrCallbackURL         string
	QrString              string
	InvoiceNumber         string
	Agent                 string
	Items                 []Item
}

type QRCallbackRequest struct {
	Event      string                          `json:"event"`
	APIVersion string                          `json:"api_version"`
	BusinessID string                          `json:"business_id"`
	Created    time.Time                       `json:"created"`
	Data       XenditWebhookRequestPaymentData `json:"data"`
}

// PaymentData struct for the 'data' field
type XenditWebhookRequestPaymentData struct {
	WebHookID     string                            `json:webhook_id`
	ID            string                            `json:"id"`
	BusinessID    string                            `json:"business_id"`
	Currency      string                            `json:"currency"`
	Amount        int                               `json:"amount"`
	Status        string                            `json:"status"`
	Created       time.Time                         `json:"created"`
	QRID          string                            `json:"qr_id"`
	QRString      string                            `json:"qr_string"`
	ReferenceID   string                            `json:"reference_id"`
	Type          string                            `json:"type"`
	ChannelCode   string                            `json:"channel_code"`
	ExpiresAt     time.Time                         `json:"expires_at"`
	Description   string                            `json:"description"`
	Basket        string                            `json:"basket"`
	Metadata      string                            `json:"metadata"`
	PaymentDetail XenditWebhookRequestPaymentDetail `json:"payment_detail"`
}

// PaymentDetail struct for the 'payment_detail' field within PaymentData
type XenditWebhookRequestPaymentDetail struct {
	ReceiptID      string `json:"receipt_id"`
	Source         string `json:"source"`
	Name           string `json:"name"`
	AccountDetails string `json:"account_details"`
}
