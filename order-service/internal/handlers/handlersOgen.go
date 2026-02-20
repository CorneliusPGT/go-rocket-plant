package handlers

import (
	"context"
	"errors"
	"fmt"
	"order-service/internal/oapi"
	api "order-service/internal/oapi"
	"order-service/internal/repository/model"
	"order-service/internal/service/order"
)

type OrderHandler struct {
	Service *order.Service
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *api.CreateOrderRequest) (api.CreateOrderRes, error) {
	if len(req.Items) == 0 {
		return nil,  fmt.Errorf("%w: items required", model.ErrBadRequest)
	}
	items := make([]model.Item, len(req.Items))
	if req.UserUUID == "" {
		return nil, model.ErrBadRequest
	}

	for i, v := range req.Items {
		if v.PartUUID == "" {
			return nil, fmt.Errorf("%w: part_uuid is required for item %d", model.ErrBadRequest, i)
		}
		if v.Quantity <= 0 {
			return nil, fmt.Errorf("%w: quantity must be greater than 0 for item %d", model.ErrBadRequest, i)
		}

		items[i] = model.Item{
			PartUUID: v.PartUUID,
			Quantity: int(v.Quantity),
		}
	}

	order, err := h.Service.CreateOrder(ctx, req.UserUUID, items)
	if err != nil {
		return nil, err
	}
	return &api.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

func (h *OrderHandler) GetOrder(
	ctx context.Context,
	params api.GetOrderParams,
) (api.GetOrderRes, error) {

	order, err := h.Service.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		return nil, err
	}
	items := make([]oapi.OrderItemsItem, 0, len(order.Items))

	for _, v := range order.Items {
		items = append(items, api.OrderItemsItem{
			Quantity: float64(v.Quantity),
			PartUUID: v.PartUUID,
			Price:    v.Price,
			Name:     v.Name,
		})
	}

	return &api.Order{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		Items:      items,
		TotalPrice: order.TotalPrice,
		Status:     api.OrderStatus(order.Status),
	}, nil
}

func (h *OrderHandler) CancelOrder(
	ctx context.Context,
	params api.CancelOrderParams,
) (api.CancelOrderRes, error) {

	err := h.Service.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		return nil, err
	}

	return &api.CancelOrderNoContent{}, nil
}

func (h *OrderHandler) PayOrder(
	ctx context.Context,
	req *api.PayOrderRequest,
	params api.PayOrderParams,
) (api.PayOrderRes, error) {

	pm := model.PaymentMethod(req.PaymentMethod)

	tUid, err := h.Service.PayOrder(ctx, params.OrderUUID, &pm)
	if err != nil {
		return nil, err
	}

	return &api.PayOrderResponse{
		TransactionUUID: tUid,
	}, nil
}

func (h *OrderHandler) NewError(
	ctx context.Context,
	err error,
) *api.ErrorStatusCode {

	switch {
	case errors.Is(err, model.ErrNotFound):
		return &api.ErrorStatusCode{
			StatusCode: 404,
			Response: api.Error{
				Message: err.Error(),
			},
		}

	case errors.Is(err, model.ErrConflict):
		return &api.ErrorStatusCode{
			StatusCode: 409,
			Response: api.Error{
				Message: err.Error(),
			},
		}

	case errors.Is(err, model.ErrBadRequest):
		return &api.ErrorStatusCode{
			StatusCode: 400,
			Response: api.Error{
				Message: err.Error(),
			},
		}

	case errors.Is(err, model.ErrNotEnoughInStock):
		return &api.ErrorStatusCode{
			StatusCode: 400,
			Response: api.Error{
				Message: err.Error(),
			},
		}

	default:
		return &api.ErrorStatusCode{
			StatusCode: 500,
			Response: api.Error{
				Message: err.Error(),
			},
		}
	}
}
