package handlers

import (
	"context"
	"errors"
	api "order-service/internal/oapi"
	"order-service/internal/repository/model"
	"order-service/internal/service/order"
)

type OrderHandler struct {
	Service *order.Service
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *api.CreateOrderRequest) (api.CreateOrderRes, error) {
	order, err := h.Service.CreateOrder(ctx, req.UserUUID, req.PartUuids)
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

	return &api.Order{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
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
				Message: "order not found",
			},
		}

	case errors.Is(err, model.ErrConflict):
		return &api.ErrorStatusCode{
			StatusCode: 409,
			Response: api.Error{
				Message: "order conflict",
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
