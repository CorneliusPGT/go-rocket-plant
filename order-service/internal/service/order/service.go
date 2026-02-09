package order

import (
	"context"
	"order-service/internal/repository"
	"order-service/internal/repository/model"
	"order-service/internal/service"

	"github.com/google/uuid"
)

type Service struct {
	repo repository.OrderRepository
	inv  service.InventoryService
	pay  service.PaymentService
}

func NewService(repo repository.OrderRepository, inv service.InventoryService, pay service.PaymentService) *Service {
	return &Service{
		repo: repo,
		inv:  inv,
		pay:  pay,
	}
}

func (s *Service) CreateOrder(ctx context.Context, userID string, partIDs []string) (*model.Order, error) {
	parts, err := s.inv.ListParts(ctx, partIDs)
	if err != nil {
		return nil, err
	}
	var total float64
	for _, part := range parts {
		total += part.Price
	}
	order := &model.Order{
		OrderUUID: uuid.New().String(),
		UserUUID:  userID,
		/* 	Parts:     parts, */
		PartUUIDs:  partIDs,
		TotalPrice: total,
		Status:     model.StatusPendingPayment,
	}

	err = s.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Service) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	return s.repo.Get(ctx, orderID)
}

func (s *Service) PayOrder(ctx context.Context, orderID string, pm *model.PaymentMethod) (string, error) {
	order, err := s.repo.Get(ctx, orderID)
	if err != nil {
		return "", err
	}
	if order.Status == model.StatusPaid {
		return "", model.ErrConflict
	}
	tId, err := s.pay.MakePayment(ctx, order.UserUUID, order.OrderUUID, pm)
	if err != nil {
		return "", err
	}
	order.PaymentMethod = pm
	order.TransactionUUID = &tId
	order.Status = model.StatusPaid
	s.repo.Update(ctx, order)
	return tId, nil
}

func (s *Service) CancelOrder(ctx context.Context, orderId string) error {
	order, err := s.repo.Get(ctx, orderId)
	if err != nil {
		return model.ErrNotFound
	}
	if order.Status == model.StatusPaid {
		return model.ErrConflict
	}
	order.Status = model.StatusCancelled
	s.repo.Update(ctx, order)
	return nil

}
