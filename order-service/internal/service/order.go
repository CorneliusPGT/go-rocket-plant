package service

import (
	"context"
	"errors"
	"inventory-service/grpc/inventorypb"
	client "order-service/cmd/http_client"
	"order-service/internal/models"

	"github.com/google/uuid"
)

type InventoryServiceImp struct {
	client *client.InventoryClient
}

type PaymentServiceImp struct {
	client *client.PaymentClient
}

func NewInventoryServiceClient(client *client.InventoryClient) *InventoryServiceImp {
	return &InventoryServiceImp{
		client: client,
	}
}

func NewPaymentServiceClient(client *client.PaymentClient) *PaymentServiceImp {
	return &PaymentServiceImp{
		client: client,
	}
}

func (s *InventoryServiceImp) GetListParts(ctx context.Context, partsIds []string) ([]*inventorypb.Part, error) {
	return s.client.ListParts(ctx, partsIds)
}

func (s *PaymentServiceImp) MakePayment(ctx context.Context, orderUuid string, userUuid string, pm *models.PaymentMethod) (string, error) {
	return s.client.MakePayment(ctx, orderUuid, userUuid, pm)
}

type OrderService struct {
	storage   *models.OrderStorage
	inventory InventoryServiceImp
	payment   PaymentServiceImp
}

func NewOrderService(
	storage *models.OrderStorage,
	inv InventoryServiceImp,
	payment PaymentServiceImp,
) *OrderService {
	return &OrderService{
		storage:   storage,
		inventory: inv,
		payment:   payment,
	}
}

var (
	ErrNotFound = errors.New("заказ не найден")
	ErrConflict = errors.New("заказ уже оплачен")
)

func (s *OrderService) CreateOrder(ctx context.Context, userID string, partIDs []string) (*models.Order, error) {

	parts, err := s.inventory.GetListParts(ctx, partIDs)
	if err != nil {
		return nil, err
	}
	var total float64
	for _, v := range parts {
		total += float64(v.Price)
	}
	order := models.Order{
		OrderUUID:  uuid.New().String(),
		UserUUID:   userID,
		PartUUIDs:  partIDs,
		TotalPrice: total,
		Status:     models.StatusPendingPayment,
	}

	err = s.storage.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) GetOrderById(ctx context.Context, orderId string) (*models.Order, error) {
	return s.storage.GetOrderById(ctx, orderId)
}

func (s *OrderService) MakePayment(ctx context.Context, pm *models.PaymentMethod, orderId string) (string, error) {
	order, err := s.storage.GetOrderById(ctx, orderId)
	if err != nil {
		return "", err
	}
	if order.Status == models.StatusPaid {
		return "", ErrConflict
	}
	tId, err := s.payment.MakePayment(ctx, order.UserUUID, order.OrderUUID, pm)
	if err != nil {
		return "", err
	}
	order.PaymentMethod = pm
	order.TransactionUUID = &tId
	order.Status = models.StatusPaid
	s.storage.UpdateOrder(ctx, order)
	return tId, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderId string) error {
	order, err := s.storage.GetOrderById(ctx, orderId)
	if err != nil {
		return ErrNotFound
	}
	if order.Status == models.StatusPaid {
		return ErrConflict
	}
	order.Status = models.StatusCancelled
	s.storage.UpdateOrder(ctx, order)
	return nil

}
