package client

import (
	"context"
	"fmt"
	inventorypb "inventory-service/grpc/inventorypb"
	"order-service/internal/models"
	paymentpb "payment-service/grpc/paymentpb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type PaymentClient struct {
	conn   *grpc.ClientConn
	client paymentpb.PaymentServiceClient
}

type InventoryClient struct {
	conn   *grpc.ClientConn
	client inventorypb.InventoryServiceClient
}

func NewInventoryClient(addr string) (*InventoryClient, error) {
	kp := keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             2 * time.Second,
		PermitWithoutStream: true,
	}
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kp),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}
	ic := &InventoryClient{
		conn:   conn,
		client: inventorypb.NewInventoryServiceClient(conn),
	}
	return ic, nil
}

func (c *InventoryClient) Close() error {
	return c.conn.Close()
}

func (c *InventoryClient) ListParts(ctx context.Context, partsIds []string) ([]*inventorypb.Part, error) {
	resp, err := c.client.ListParts(ctx, &inventorypb.ListPartsRequest{
		Filter: &inventorypb.PartsFilter{
			Uuids: partsIds,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.Parts, nil
}

func NewPaymentClient(addr string) (*PaymentClient, error) {
	kp := keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             2 * time.Second,
		PermitWithoutStream: true,
	}
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kp),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}
	pc := &PaymentClient{
		conn:   conn,
		client: paymentpb.NewPaymentServiceClient(conn),
	}

	return pc, nil
}

func (c *PaymentClient) Close() error {
	return c.conn.Close()
}

func (c *PaymentClient) MakePayment(ctx context.Context, orderUuid, userUuid string, pm *models.PaymentMethod) (string, error) {
	if pm == nil {
		return "", fmt.Errorf("payment method is required")
	}
	pmProto := paymentpb.PaymentMethod_value[string(*pm)]

	req := &paymentpb.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: paymentpb.PaymentMethod(pmProto),
	}
	resp, err := c.client.PayOrder(ctx, req)
	if err != nil {
		return "", fmt.Errorf("PayOrder failed: %w", err)
	}

	return resp.TransactionUuid, nil
}
