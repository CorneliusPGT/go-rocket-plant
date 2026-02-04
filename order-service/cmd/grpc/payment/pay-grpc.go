package paymentgrpc

import (
	"context"
	"fmt"
	"order-service/internal/repository/model"
	"payment-service/grpc/paymentpb"
)

type GRPCClient struct {
	client paymentpb.PaymentServiceClient
}

func New(client paymentpb.PaymentServiceClient) *GRPCClient {
	return &GRPCClient{
		client: client,
	}
}

func (g *GRPCClient) MakePayment(ctx context.Context, orderID, userID string, pm *model.PaymentMethod) (string, error) {
	if pm == nil {
		return "", fmt.Errorf("payment method is required")
	}
	pbValue, ok := paymentpb.PaymentMethod_value[string(*pm)]
	if !ok {
		pbValue = int32(paymentpb.PaymentMethod_UNKNOWN)
	}
	resp, err := g.client.PayOrder(ctx, &paymentpb.PayOrderRequest{
		OrderUuid:     orderID,
		UserUuid:      userID,
		PaymentMethod: paymentpb.PaymentMethod(pbValue),
	})
	if err != nil {
		return "", err
	}
	return resp.TransactionUuid, nil
}
