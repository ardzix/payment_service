package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payment-service/internal/domain"
	"payment-service/internal/usecase"
)

type PaymentHandler struct {
	useCase usecase.PaymentUseCase
}

func NewPaymentHandler(useCase usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{useCase: useCase}
}

// Example handler for creating a payment
func (c *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var payload domain.QRCallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var webhookId = r.Header.Get("webhook-id")
	fmt.Println(&payload)
	fmt.Println(webhookId)

	data, err := c.useCase.QrWebhook(r.Context(), payload.Data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)

}
