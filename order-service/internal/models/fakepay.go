package models

import "context"

type PaymentService interface {
	MakePayment(context.Context, string, string, *PaymentMethod) (*string, error)
}
