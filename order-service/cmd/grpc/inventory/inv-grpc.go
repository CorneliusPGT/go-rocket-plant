package inventorygrpc

import (
	"context"
	"inventory-service/grpc/inventorypb"
	"order-service/internal/repository/model"
)

type GRPCClient struct {
	client inventorypb.InventoryServiceClient
}

func New(client inventorypb.InventoryServiceClient) *GRPCClient {
	return &GRPCClient{client: client}
}

func (g *GRPCClient) ListParts(ctx context.Context, partIDs []string) ([]*model.Part, error) {
	resp, err := g.client.ListParts(ctx, &inventorypb.ListPartsRequest{
		Filter: &inventorypb.PartsFilter{
			Uuids: partIDs,
		},
	})
	if err != nil {
		return nil, err
	}
	parts := make([]*model.Part, len(resp.Parts))
	for i, v := range resp.Parts {
		parts[i] = &model.Part{
			Quantity: int(v.StockQuantity),
			UUID:     v.Uuid,
			Name:     v.Name,
			Price:    v.Price,
		}
	}
	return parts, nil
}
