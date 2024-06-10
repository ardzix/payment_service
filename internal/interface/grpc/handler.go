// internal/interface/grpc/handler.go
package grpc

import (
	"context"
	"errors"
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

	// Validate amount matches total of items' price * quantity
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}
	if totalAmount != req.Amount {
		return nil, errors.New("amount does not match the total of items' price * quantity")
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
		return nil, err
	}

	return &proto.ProcessPaymentResponse{
		PaymentId: result.PaymentID,
		Status:    result.Status,
	}, nil
}

func (h *PaymentHandler) RefundPayment(ctx context.Context, req *proto.RefundPaymentRequest) (*proto.RefundPaymentResponse, error) {
	refundID, err := h.useCase.RefundPayment(ctx, req.PaymentId, req.Amount)
	if err != nil {
		return nil, err
	}

	return &proto.RefundPaymentResponse{RefundId: refundID, Status: "refunded"}, nil
}

func (h *PaymentHandler) GetPaymentStatus(ctx context.Context, req *proto.GetPaymentStatusRequest) (*proto.GetPaymentStatusResponse, error) {
	payment, err := h.useCase.GetPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, err
	}

	return &proto.GetPaymentStatusResponse{
		PaymentId:    payment.PaymentID,
		Status:       payment.Status,
		ErrorMessage: "",
	}, nil
}

func (h *PaymentHandler) ListPayments(ctx context.Context, req *proto.ListPaymentsRequest) (*proto.ListPaymentsResponse, error) {
	payments, total, err := h.useCase.ListPayments(ctx, req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
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

	return response, nil
}

func (h *PaymentHandler) GetPaymentDetail(ctx context.Context, req *proto.GetPaymentDetailRequest) (*proto.GetPaymentDetailResponse, error) {
	payment, err := h.useCase.GetPayment(ctx, req.PaymentId)
	if err != nil {
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
		InvoiceNumber:         payment.InvoiceNumber,
		Agent:                 payment.Agent,
		Items:                 items,
	}, nil
}
