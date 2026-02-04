package service

import (
	"context"
	"order-service/internal/repository/model"
)

type InventoryService interface {
	ListParts(ctx context.Context, partIDs []string) ([]*model.Part, error)
}

type PaymentService interface {
	MakePayment(ctx context.Context, orderID, userID string, pm *model.PaymentMethod) (string, error)
}

type OrderService interface {
	CreateOrder(ctx context.Context, userID string, partIDs []string) (*model.Order, error)
	GetOrder(ctx context.Context, orderID string) (*model.Order, error)
	PayOrder(ctx context.Context, orderID string, pm *model.PaymentMethod) (string, error)
	CancelOrder(ctx context.Context, orderID string) error
}
