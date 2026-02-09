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
	ErrConflict = errors.New("409 conflict")
	ErrNotFound = errors.New("404 not found")
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
	 PartUUIDs       []string `json:"part_uuids"` 
/* 	Parts                          []*Part `json:"parts"` */
	TotalPrice                     float64 `json:"total_price"`
	TransactionUUID                *string
	PaymentMethod                  *PaymentMethod `json:"payment_method"`
	Status                         OrderStatus    `json:"status"`
}

type Part struct {
	UUID     string
	Name     string
	Price    float64
	Quantity int
}
