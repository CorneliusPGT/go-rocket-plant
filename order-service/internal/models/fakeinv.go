package models

import "context"

type InventoryService interface {
	GetListParts(context.Context, string, []string) (float64, error)
}
