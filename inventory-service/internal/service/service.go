package service

import (
	"context"
	"inventory-service/grpc/inventorypb"
	repo "inventory-service/repository"
)

type PartService interface {
	Get(ctx context.Context, uuid string) (*inventorypb.Part, error)
	List(ctx context.Context, filter *inventorypb.PartsFilter) ([]*inventorypb.Part, error)
}

type Service struct {
	repo repo.PartRepo
}

func NewPartService(r repo.PartRepo) PartService {
	return &Service{repo: r}
}

func (s *Service) Get(ctx context.Context, uuid string) (*inventorypb.Part, error) {
	return s.repo.Get(ctx, uuid)
}

func (s *Service) List(ctx context.Context, filter *inventorypb.PartsFilter) ([]*inventorypb.Part, error) {
	return s.repo.List(ctx, filter)
}
