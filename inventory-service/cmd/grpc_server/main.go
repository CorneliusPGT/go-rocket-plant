package main

import (
	"context"
	"fmt"
	"inventory-service/grpc/inventorypb"

	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type inventoryService struct {
	inventorypb.UnimplementedInventoryServiceServer
	mu    sync.RWMutex
	parts map[string]*inventorypb.Part
}

func seedParts() map[string]*inventorypb.Part {
	now := timestamppb.Now()

	return map[string]*inventorypb.Part{
		"engine-1": {
			Uuid:          "engine-1",
			Name:          "Main Engine",
			Description:   "Primary propulsion engine",
			Price:         1_500_000,
			StockQuantity: 10,
			Category:      inventorypb.Category_CATEGORY_ENGINE,
			Dimensions: &inventorypb.Dimensions{
				Length: 4,
				Width:  2,
				Height: 2,
				Weight: 1500,
			},
			Manufacter: &inventorypb.Manufacter{
				Name:    "SpaceY",
				Country: "USA",
				Website: "https://spacey.example",
			},
			Tags:      []string{"engine", "rocket"},
			CreatedAt: now,
			UpdatedAt: now,
		},

		"wing-1": {
			Uuid:          "wing-1",
			Name:          "Left Wing",
			Description:   "Aerodynamic wing",
			Price:         250_000,
			StockQuantity: 5,
			Category:      inventorypb.Category_CATEGORY_WING,
			Manufacter: &inventorypb.Manufacter{
				Name:    "AeroWorks",
				Country: "Germany",
			},
			Tags:      []string{"wing"},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()
	service := &inventoryService{
		parts: seedParts(),
	}
	inventorypb.RegisterInventoryServiceServer(s, service)
	reflection.Register(s)

	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server stopped")
}

func (s *inventoryService) GetPart(ctx context.Context, req *inventorypb.GetPartRequest) (*inventorypb.GetPartResponse, error) {
	s.mu.RLock()
	part, ok := s.parts[req.Uuid]
	s.mu.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "part not found")
	}
	return &inventorypb.GetPartResponse{Part: part}, nil
}

func (s *inventoryService) ListParts(ctx context.Context, req *inventorypb.ListPartsRequest) (*inventorypb.ListPartsResponse, error) {
	var parts []*inventorypb.Part
	s.mu.RLock()
	for _, v := range s.parts {
		parts = append(parts, v)
	}
	s.mu.RUnlock()
	if req.Filter == nil {
		return &inventorypb.ListPartsResponse{Parts: parts}, nil
	}
	var filtered []*inventorypb.Part
	for _, part := range parts {
		if checkPart(part, req.Filter) {
			filtered = append(filtered, part)
		}
	}
	return &inventorypb.ListPartsResponse{Parts: filtered}, nil
}

func checkPart(part *inventorypb.Part, filters *inventorypb.PartsFilter) bool {
	fCheck := false
	if len(filters.Uuids) > 0 {
		for _, v := range filters.Uuids {
			if part.Uuid == v {
				fCheck = true
			}
		}
		if fCheck == false {
			return false
		}
		fCheck = false
	}
	if len(filters.Names) > 0 {
		for _, v := range filters.Names {
			if part.Name == v {
				fCheck = true
			}
		}
		if fCheck == false {
			return false
		}
		fCheck = false
	}
	if len(filters.Categories) > 0 {
		for _, v := range filters.Categories {
			if part.Category == v {
				fCheck = true
			}
		}
		if fCheck == false {
			return false
		}
		fCheck = false
	}
	if len(filters.ManufacturerCountries) > 0 {
		for _, v := range filters.ManufacturerCountries {
			if part.Manufacter.Country == v {
				fCheck = true
			}
		}
		if fCheck == false {
			return false
		}
		fCheck = false
	}
	if len(filters.Tags) > 0 {
		for _, v := range filters.Tags {
			for _, k := range part.Tags {
				if k == v {
					fCheck = true
				}
			}
		}
		if fCheck == false {
			return false
		}

	}
	return true
}
