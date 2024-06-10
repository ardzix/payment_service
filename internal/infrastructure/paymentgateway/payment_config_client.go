package paymentgateway

import (
	"context"
	"payment-service/api/proto"
	"time"

	"google.golang.org/grpc"
)

type PaymentConfigClient struct {
	client  proto.PaymentConfigServiceClient
	timeout time.Duration
}

func NewPaymentConfigClient(conn *grpc.ClientConn, timeout time.Duration) *PaymentConfigClient {
	return &PaymentConfigClient{
		client:  proto.NewPaymentConfigServiceClient(conn),
		timeout: timeout,
	}
}

func (pcc *PaymentConfigClient) GetPaymentGatewayConfig(paymentMethod string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pcc.timeout)
	defer cancel()

	req := &proto.PaymentMethodRequest{
		PaymentMethod: paymentMethod,
	}

	resp, err := pcc.client.GetPaymentGatewayConfig(ctx, req)
	if err != nil {
		return "", "", err
	}

	return resp.Gateway, resp.FallbackGateway, nil
}
