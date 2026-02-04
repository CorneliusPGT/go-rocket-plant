package order

import (
	"context"
	"errors"
	"order-service/internal/mocks"
	"order-service/internal/repository/model"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OrderServiceTest struct {
	suite.Suite

	repo *mocks.OrderRepository
	inv  *mocks.InventoryService
	pay  *mocks.PaymentService

	service *Service
}

func (s *OrderServiceTest) SetupTest() {
	s.repo = mocks.NewOrderRepository(s.T())
	s.inv = mocks.NewInventoryService(s.T())
	s.pay = mocks.NewPaymentService(s.T())

	s.service = NewService(s.repo, s.inv, s.pay)
}

func TestOrderServiceTest(t *testing.T) {
	suite.Run(t, new(OrderServiceTest))
}

func (s *OrderServiceTest) TestCreateOrder_success() {
	ctx := context.Background()

	partIDs := []string{"p1", "p2"}

	s.inv.On("ListParts", ctx, partIDs).Return([]*model.Part{
		{UUID: "p1", Price: 10},
		{UUID: "p2", Price: 20},
	}, nil)
	s.repo.On("Create", ctx, mock.AnythingOfType("*model.Order")).Return(nil)
	order, err := s.service.CreateOrder(ctx, "user-1", partIDs)

	s.NoError(err)
	s.Equal(float64(30), order.TotalPrice)

	s.inv.AssertExpectations(s.T())
	s.repo.AssertExpectations(s.T())
}

func (s *OrderServiceTest) TestCreateOrder_inventoryError() {
	ctx := context.Background()

	partIDs := []string{"p1", "p2"}

	s.inv.On("ListParts", ctx, partIDs).Return(nil, errors.New("not found"))
	_, err := s.service.CreateOrder(ctx, "user-1", partIDs)
	s.Error(err)
	s.inv.AssertExpectations(s.T())
	s.repo.AssertNotCalled(s.T(), "Create", mock.Anything)
}
func (s *OrderServiceTest) TestPayOrder_success() {
	ctx := context.Background()

	orderID := "id-1"
	userID := "u-1"
	order := &model.Order{
		OrderUUID: orderID,
		UserUUID:  userID,
		Status:    model.StatusPendingPayment,
	}
	s.repo.On("Get", ctx, orderID).Return(order, nil)
	s.pay.On("MakePayment", ctx, userID, orderID, (*model.PaymentMethod)(nil)).Return("tId-1", nil)
	s.repo.On("Update", ctx, mock.AnythingOfType("*model.Order")).Return(nil)
	_, err := s.service.PayOrder(ctx, orderID, nil)
	s.NoError(err)
	s.repo.AssertExpectations(s.T())
	s.pay.AssertExpectations(s.T())
}
func (s *OrderServiceTest) TestCancelOrder_conflict() {
	ctx := context.Background()

	orderId := "id-1"
	s.repo.On("Get", ctx, orderId).Return(&model.Order{
		OrderUUID: orderId,
		Status:    model.StatusPaid,
	}, nil)
	err := s.service.CancelOrder(ctx, orderId)
	s.Error(err)
	s.repo.AssertNotCalled(s.T(), "Update", mock.Anything)
}
