// internal/interface/grpc/handler.go
package grpc

import (
	"context"
	"errors"
	"log"
	"payment-service/api/proto"
	"payment-service/internal/domain"
	"payment-service/internal/usecase"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentHandler struct {
	proto.UnimplementedPaymentServiceServer
	useCase usecase.PaymentUseCase
}

func NewPaymentHandler(useCase usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{useCase: useCase}
}

func (h *PaymentHandler) ProcessPayment(ctx context.Context, req *proto.ProcessPaymentRequest) (*proto.ProcessPaymentResponse, error) {
	log.Printf("Received ProcessPayment request: UserId=%s, Amount=%.2f, InvoiceNumber=%s", req.UserId, req.Amount, req.InvoiceNumber)

	// Validate required fields
	if req.UserId == "" {
		return nil, errors.New("user_id is required")
	}
	if req.Amount == 0 {
		return nil, errors.New("amount is required")
	}
	if req.InvoiceNumber == "" {
		return nil, errors.New("invoice_number is required")
	}
	if req.Agent == "" {
		return nil, errors.New("agent is required")
	}
	if len(req.Items) == 0 {
		return nil, errors.New("items are required")
	}

	items := make([]domain.Item, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.Item{
			ItemName: item.ItemName,
			Quantity: int(item.Quantity),
			Price:    item.Price,
		}
	}

	payment := &domain.Payment{
		PaymentID:             uuid.New().String(),
		UserID:                req.UserId,
		Amount:                req.Amount,
		Currency:              req.Currency,
		PaymentMethod:         req.PaymentMethod,
		PhoneNumber:           req.PhoneNumber,
		EwalletCheckoutMethod: req.EwalletCheckoutMethod,
		QrType:                req.QrType,
		QrCallbackURL:         req.QrCallbackUrl,
		InvoiceNumber:         req.InvoiceNumber,
		Agent:                 req.Agent,
		Items:                 items,
	}

	result, err := h.useCase.ProcessPayment(ctx, payment)
	if err != nil {
		log.Printf("Error processing payment: %v", err)
		return nil, err
	}

	log.Printf("Payment processed successfully: PaymentId=%s, Status=%s", result.PaymentID, result.Status)
	return &proto.ProcessPaymentResponse{
		PaymentId:     result.PaymentID,
		Status:        result.Status,
		PaymentMethod: result.PaymentMethod,
		QrString:      result.QrString,
	}, nil
}

func (h *PaymentHandler) RefundPayment(ctx context.Context, req *proto.RefundPaymentRequest) (*proto.RefundPaymentResponse, error) {
	log.Printf("Received RefundPayment request: PaymentId=%s, Amount=%.2f", req.PaymentId, req.Amount)

	refundID, err := h.useCase.RefundPayment(ctx, req.PaymentId, req.Amount)
	if err != nil {
		log.Printf("Error refunding payment: %v", err)
		return nil, err
	}

	log.Printf("Payment refunded successfully: RefundId=%s", refundID)

	return &proto.RefundPaymentResponse{RefundId: refundID, Status: "refunded"}, nil
}

func (h *PaymentHandler) GetPaymentStatus(ctx context.Context, req *proto.GetPaymentStatusRequest) (*proto.GetPaymentStatusResponse, error) {
	log.Printf("Received GetPaymentStatus request: PaymentId=%s", req.PaymentId)

	payment, err := h.useCase.GetPayment(ctx, req.PaymentId)
	if err != nil {
		log.Printf("Error getting payment status: %v", err)
		return nil, err
	}

	log.Printf("Payment status retrieved successfully: PaymentId=%s, Status=%s", payment.PaymentID, payment.Status)

	return &proto.GetPaymentStatusResponse{
		PaymentId:    payment.PaymentID,
		Status:       payment.Status,
		ErrorMessage: "",
	}, nil
}

func (h *PaymentHandler) ListPayments(ctx context.Context, req *proto.ListPaymentsRequest) (*proto.ListPaymentsResponse, error) {
	log.Printf("Received ListPayments request: UserId=%s, Page=%d, PageSize=%d", req.UserId, req.Page, req.PageSize)

	payments, total, err := h.useCase.ListPayments(ctx, req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
		log.Printf("Error listing payments: %v", err)
		return nil, err
	}

	response := &proto.ListPaymentsResponse{TotalCount: int32(total)}
	for _, payment := range payments {
		response.Payments = append(response.Payments, &proto.Payment{
			PaymentId:     payment.PaymentID,
			UserId:        payment.UserID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			PaymentMethod: payment.PaymentMethod,
			Gateway:       payment.Gateway,
			Status:        payment.Status,
			CreatedAt:     timestamppb.New(payment.CreatedAt),
		})
	}

	log.Printf("Payments listed successfully for UserId=%s", req.UserId)

	return response, nil
}

func (h *PaymentHandler) GetPaymentDetail(ctx context.Context, req *proto.GetPaymentDetailRequest) (*proto.GetPaymentDetailResponse, error) {
	log.Printf("Received GetPaymentDetail request: PaymentId=%s", req.PaymentId)

	payment, err := h.useCase.GetPayment(ctx, req.PaymentId)
	if err != nil {
		log.Printf("Error getting payment detail: %v", err)
		return nil, err
	}

	items := make([]*proto.Item, len(payment.Items))
	for i, item := range payment.Items {
		items[i] = &proto.Item{
			ItemName: item.ItemName,
			Quantity: int32(item.Quantity),
			Price:    item.Price,
		}
	}

	log.Printf("Payment detail retrieved successfully: PaymentId=%s", payment.PaymentID)
	return &proto.GetPaymentDetailResponse{
		PaymentId:             payment.PaymentID,
		UserId:                payment.UserID,
		Amount:                payment.Amount,
		Currency:              payment.Currency,
		Status:                payment.Status,
		CreatedAt:             timestamppb.New(payment.CreatedAt),
		UpdatedAt:             timestamppb.New(payment.UpdatedAt),
		PaymentMethod:         payment.PaymentMethod,
		Gateway:               payment.Gateway,
		PhoneNumber:           payment.PhoneNumber,
		EwalletCheckoutMethod: payment.EwalletCheckoutMethod,
		QrType:                payment.QrType,
		QrCallbackUrl:         payment.QrCallbackURL,
		QrString:              payment.QrString,
		InvoiceNumber:         payment.InvoiceNumber,
		Agent:                 payment.Agent,
		Items:                 items,
	}, nil
}
