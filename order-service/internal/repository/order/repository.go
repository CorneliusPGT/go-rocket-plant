package repository

import (
	"context"
	"errors"
	"order-service/internal/repository/model"
	"sync"
)

type Repository struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
}

func NewRepository() *Repository {
	return &Repository{
		orders: make(map[string]*model.Order),
	}
}

func (o *Repository) Create(ctx context.Context, order *model.Order) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if _, exists := o.orders[order.OrderUUID]; exists {
		return errors.New("заказ с таким UUID уже существует")
	}
	o.orders[order.OrderUUID] = order
	return nil
}

func (o *Repository) Get(ctx context.Context, orderId string) (*model.Order, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	order, exists := o.orders[orderId]
	if !exists {
		return nil, errors.New("заказ не найден")
	}
	return order, nil
}

func (o *Repository) Update(ctx context.Context, order *model.Order) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.orders[order.OrderUUID] = order
	return nil
}
