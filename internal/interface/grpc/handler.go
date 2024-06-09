// internal/interface/grpc/handler.go
package grpc

import (
	"context"
	"payment-service/api/proto"
	"payment-service/internal/domain"
	"payment-service/internal/usecase"

	"github.com/google/uuid"
)

type PaymentHandler struct {
	proto.UnimplementedPaymentServiceServer
	useCase usecase.PaymentUseCase
}

func NewPaymentHandler(useCase usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{useCase: useCase}
}

func (h *PaymentHandler) ProcessPayment(ctx context.Context, req *proto.ProcessPaymentRequest) (*proto.ProcessPaymentResponse, error) {
	payment := &domain.Payment{
		PaymentID:     uuid.New().String(),
		UserID:        req.UserId,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
	}

	result, err := h.useCase.ProcessPayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	return &proto.ProcessPaymentResponse{PaymentId: result.PaymentID, Status: result.Status}, nil
}

func (h *PaymentHandler) RefundPayment(ctx context.Context, req *proto.RefundPaymentRequest) (*proto.RefundPaymentResponse, error) {
	refundID, err := h.useCase.RefundPayment(ctx, req.PaymentId, req.Amount)
	if err != nil {
		return nil, err
	}

	return &proto.RefundPaymentResponse{RefundId: refundID, Status: "refunded"}, nil
}

func (h *PaymentHandler) GetPaymentStatus(ctx context.Context, req *proto.GetPaymentStatusRequest) (*proto.GetPaymentStatusResponse, error) {
	payment, err := h.useCase.GetPaymentStatus(ctx, req.PaymentId)
	if err != nil {
		return nil, err
	}

	return &proto.GetPaymentStatusResponse{PaymentId: payment.PaymentID, Status: payment.Status}, nil
}

func (h *PaymentHandler) ListPayments(ctx context.Context, req *proto.ListPaymentsRequest) (*proto.ListPaymentsResponse, error) {
	payments, totalCount, err := h.useCase.ListPayments(ctx, req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	response := &proto.ListPaymentsResponse{TotalCount: int32(totalCount)}
	for _, payment := range payments {
		response.Payments = append(response.Payments, &proto.Payment{
			PaymentId: payment.PaymentID,
			UserId:    payment.UserID,
			Amount:    payment.Amount,
			Currency:  payment.Currency,
			Status:    payment.Status,
			CreatedAt: payment.CreatedAt,
		})
	}

	return response, nil
}
