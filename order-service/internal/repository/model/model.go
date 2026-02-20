package model

import (
	"errors"
)

type OrderStatus string

type PaymentMethod string

const (
	StatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	StatusPaid           OrderStatus = "PAID"
	StatusCancelled      OrderStatus = "CANCELLED"
)

var (
	ErrBadRequest       = errors.New("400 bad request")
	ErrConflict         = errors.New("409 conflict")
	ErrNotFound         = errors.New("404 not found")
	ErrNotEnoughInStock = errors.New("400 not enough in stock")
)

const (
	PaymentCard     PaymentMethod = "CARD"
	PaymentSBP      PaymentMethod = "SBP"
	PaymentCredit   PaymentMethod = "CREDIT_CARD"
	PaymentInvestor PaymentMethod = "INVESTOR_MONEY"
)

type Order struct {
	OrderUUID string `json:"order_uuid"`
	UserUUID  string /* `json:"user_uuid"` */
	/* PartUUIDs       []string `json:"part_uuids"` */
	Items           []Item  `json:"items"`
	TotalPrice      float64 `json:"total_price"`
	TransactionUUID *string
	PaymentMethod   *PaymentMethod `json:"payment_method"`
	Status          OrderStatus    `json:"status"`
}

type Part struct {
	UUID     string
	Price    float64
	Quantity int
	Name     string
}

type Item struct {
	PartUUID string  `json:"part_uuid"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Name     string  `json:"name"`
}
