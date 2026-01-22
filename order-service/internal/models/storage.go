package models

import (
	"context"
	"errors"
	"sync"
)

type OrderStatus string

type PaymentMethod string


const (
	StatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	StatusPaid           OrderStatus = "PAID"
	StatusCancelled      OrderStatus = "CANCELLED"
)

const (
	PaymentCard     PaymentMethod = "CARD"
	PaymentSBP      PaymentMethod = "SBP"
	PaymentCredit   PaymentMethod = "CREDIT_CARD"
	PaymentInvestor PaymentMethod = "INVESTOR_MONEY"
)

const ()

type Order struct {
	OrderUUID       string   `json:"order_uuid"`
	UserUUID        string   `json:"user_uuid"`
	PartUUIDs       []string `json:"part_uuids"`
	TotalPrice      float64  `json:"total_price"`
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
}

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*Order),
	}
}

func (o *OrderStorage) CreateOrder(ctx context.Context, order Order) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if _, exists := o.orders[order.OrderUUID]; exists {
		return errors.New("заказ с таким UUID уже существует")
	}
	o.orders[order.OrderUUID] = &order
	return nil
}

func (o *OrderStorage) GetOrderById(ctx context.Context, orderId string) (*Order, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	order, exists := o.orders[orderId]
	if !exists {
		return nil, errors.New("заказ не найден")
	}
	return order, nil
}

func (o *OrderStorage) UpdateOrder(ctx context.Context, order *Order) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.orders[order.OrderUUID] = order
	return nil
}
