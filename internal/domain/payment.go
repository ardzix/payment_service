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
	InvoiceNumber         string
	Agent                 string
	Items                 []Item
}
