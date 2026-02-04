package repository

import (
	"context"
	"order-service/internal/repository/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	Get(ctx context.Context, orderID string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}
