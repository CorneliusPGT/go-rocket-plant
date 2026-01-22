package handlers

import (
	"encoding/json"
	"net/http"
	"order-service/internal/models"
	"order-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ResponseType struct {
}

func PostOrders(s *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID  string   `json:"user_uuid"`
			PartIDs []string `json:"part_uuids"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "неверный JSON", http.StatusBadRequest)
			return
		}

		order, err := s.CreateOrder(r.Context(), req.UserID, req.PartIDs)
		if err != nil {
			http.Error(w, "не удалось создать заказ", http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, struct {
			OrderID    string  `json:"order_uuid"`
			TotalPrice float64 `json:"total_price"`
		}{
			OrderID:    order.OrderUUID,
			TotalPrice: order.TotalPrice,
		})
	}
}

func GetOrderById(s *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderId := chi.URLParam(r, "order_uuid")

		if orderId == "" {
			http.Error(w, "order_uuid не может быть пуст", http.StatusBadRequest)
			return
		}
		order, err := s.GetOrderById(r.Context(), orderId)
		if err != nil {
			http.Error(w, "не удалось получить заказ: "+err.Error(), http.StatusNotFound)
			return
		}
		render.JSON(w, r, order)
	}
}

func MakePayment(s *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderId := chi.URLParam(r, "order_uuid")
		var req struct {
			PaymentMethod models.PaymentMethod `json:"payment_method"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "неверный JSON", http.StatusBadRequest)
			return
		}
		tId, err := s.MakePayment(r.Context(), &req.PaymentMethod, orderId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, struct {
			TransactionId *string `json:"transaction_uuid"`
		}{
			TransactionId: &tId,
		})

	}
}

func CancelOrder(s *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderId := chi.URLParam(r, "order_uuid")
		if orderId == "" {
			http.Error(w, "order_uuid не может быть пуст", http.StatusBadRequest)
			return
		}
		err := s.CancelOrder(r.Context(), orderId)
		if err.Error() == "заказ не найден" {
			http.Error(w, "заказ не найден", http.StatusNotFound)
			return
		} else if err.Error() == "заказ уже оплачен" {
			http.Error(w, "заказ уже оплачен и не может быть отменён", http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	}
}
