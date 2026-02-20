package handlers

import (
	"context"
	"errors"
	"inventory-service/grpc/inventorypb"
	"inventory-service/internal/service"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryHandler struct {
	inventorypb.UnimplementedInventoryServiceServer
	service service.PartService
}

func NewInventoryHandler(s service.PartService) *InventoryHandler {
	return &InventoryHandler{
		service: s,
	}
}

func (h *InventoryHandler) GetPart(ctx context.Context, req *inventorypb.GetPartRequest) (*inventorypb.GetPartResponse, error) {
	if req == nil || req.Uuid == "" {
		return nil, status.Errorf(codes.InvalidArgument, "uuid is required")
	}
	part, err := h.service.Get(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "part not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}
	return &inventorypb.GetPartResponse{
		Part: part,
	}, nil
}

func (h *InventoryHandler) ListParts(ctx context.Context, req *inventorypb.ListPartsRequest) (*inventorypb.ListPartsResponse, error) {
	parts, err := h.service.List(ctx, req.GetFilter())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &inventorypb.ListPartsResponse{
		Parts: parts,
	}, nil
}
