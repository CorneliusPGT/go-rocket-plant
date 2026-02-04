package main

import (
	"context"
	"inventory-service/grpc/inventorypb"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newTestService() *inventoryService {
	return &inventoryService{
		parts: seedParts(),
	}
}

func TestGetPart_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	req := &inventorypb.GetPartRequest{Uuid: "engine-1"}
	resp, err := svc.GetPart(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Part.Uuid != "engine-1" {
		t.Errorf("expected part UUID 'engine-1', got '%s'", resp.Part.Uuid)
	}
}

func TestGetPart_NotFound(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	req := &inventorypb.GetPartRequest{Uuid: "Naeblya"}
	_, err := svc.GetPart(ctx, req)
	if err == nil {
		t.Fatalf("expected error")
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.NotFound {
		t.Errorf("expected not found, gpt %v", err)
	}
}
