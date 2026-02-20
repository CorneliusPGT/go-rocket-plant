package e2e

import (
	"context"
	"order-service/internal/oapi"
	"order-service/internal/repository/model"

	"github.com/stretchr/testify/mock"
)

func (s *OrderE2ESuite) TestCreate_Success() {
	ctx := context.Background()
	req := []string{
		"engine-1",
	}

	s.Env.InvMock.On("ListParts", mock.Anything, req).Return([]*model.Part{
		{
			UUID:     "engine-1",
			Name:     "Engine",
			Price:    100,
			Quantity: 10,
		},
	}, nil).Once()

	resp, err := s.Client.CreateOrder(ctx, &oapi.CreateOrderRequest{
		UserUUID: "user-1",
		Items: []oapi.CreateOrderRequestItemsItem{
			{PartUUID: "engine-1",
				Quantity: 5},
		},
	})
	s.Require().NoError(err)
	createResp, ok := resp.(*oapi.CreateOrderResponse)
	s.Require().True(ok)

	s.NotEmpty(createResp.OrderUUID)
	s.Equal(500.0, createResp.TotalPrice)

	var count int
	var oCount int
	s.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders").Scan(&count)
	s.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM order_items").Scan(&oCount)
	s.Equal(1, count)
	s.Equal(len(req), oCount)
}

func (s *OrderE2ESuite) TestCreate_NotFound() {
	ctx := context.Background()
	req := []string{
		"engine1",
	}

	s.Env.InvMock.On("ListParts", mock.Anything, req).Return([]*model.Part{}, nil).Once()

	_, err := s.Client.CreateOrder(ctx, &oapi.CreateOrderRequest{
		UserUUID: "1",
		Items: []oapi.CreateOrderRequestItemsItem{
			{
				PartUUID: "engine1",
				Quantity: 5,
			},
		},
	})
	s.Require().Error(err)
	var count int
	s.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders").Scan(&count)
	s.Equal(0, count)
}

func (s *OrderE2ESuite) TestCreate_Zero() {
	ctx := context.Background()

	resp, err := s.Client.CreateOrder(ctx, &oapi.CreateOrderRequest{
		UserUUID: "1",
		Items: []oapi.CreateOrderRequestItemsItem{
			{
				PartUUID: "engine-1",
				Quantity: 0,
			},
		},
	})

	s.Require().NoError(err)

	_, ok := resp.(*oapi.CreateOrderBadRequest)
	s.Require().True(ok)

	var count int
	s.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders").Scan(&count)
	s.Equal(0, count)

	s.Env.InvMock.AssertNotCalled(s.T(), "ListParts", mock.Anything, mock.Anything)
}
func (s *OrderE2ESuite) TestCreate_OutOfStock() {
	ctx := context.Background()
	req := []string{"engine-1"}

	s.Env.InvMock.On("ListParts", mock.Anything, req).Return(nil, model.ErrNotEnoughInStock)

	resp, err := s.Client.CreateOrder(ctx, &oapi.CreateOrderRequest{
		UserUUID: "1",
		Items: []oapi.CreateOrderRequestItemsItem{
			{PartUUID: "engine-1", Quantity: 11},
		},
	})
	s.Require().NoError(err)
	badReq, ok := resp.(*oapi.CreateOrderBadRequest)
	s.Require().True(ok)
	s.Equal("400 not enough in stock", badReq.Message)
}

func (s *OrderE2ESuite) TestCreate_FewSuccess() {
	ctx := context.Background()
	req := []string{"engine-1", "wing-1"}
	s.Env.InvMock.On("ListParts", mock.Anything, req).Return([]*model.Part{
		{
			UUID:     "engine-1",
			Name:     "Engine",
			Price:    100,
			Quantity: 10,
		},
		{
			UUID:     "wing-1",
			Name:     "Wing",
			Price:    200,
			Quantity: 5,
		},
	}, nil)
	resp, err := s.Client.CreateOrder(ctx, &oapi.CreateOrderRequest{
		UserUUID: "1",
		Items: []oapi.CreateOrderRequestItemsItem{
			{
				PartUUID: "engine-1",
				Quantity: 8,
			},
			{
				PartUUID: "wing-1",
				Quantity: 3,
			},
		},
	})
	s.Require().NoError(err)
	createResp, ok := resp.(*oapi.CreateOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(1400.0, createResp.TotalPrice)
	s.Env.InvMock.AssertExpectations(s.T())
}

func (s *OrderE2ESuite) TestCreate_EmptyUserUUID() {
	ctx := context.Background()

	resp, _ := s.Client.CreateOrder(ctx, &oapi.CreateOrderRequest{
		UserUUID: "",
		Items: []oapi.CreateOrderRequestItemsItem{
			{
				PartUUID: "engine-1",
				Quantity: 5,
			},
		},
	})
	_, ok := resp.(*oapi.CreateOrderBadRequest)
	s.Require().True(ok)
}

func (s *OrderE2ESuite) TestCreate_EmptyItems() {
	ctx := context.Background()

	resp, err := s.Client.CreateOrder(ctx, &oapi.CreateOrderRequest{
		UserUUID: "0",
		Items:    []oapi.CreateOrderRequestItemsItem{},
	})
	_, ok := resp.(*oapi.CreateOrderBadRequest)
	s.Require().NoError(err)
	s.Require().True(ok)
}
