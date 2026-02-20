package service

import (
	"context"
	"errors"
	"testing"

	"inventory-service/grpc/inventorypb"
	"inventory-service/mocks"

	"github.com/stretchr/testify/suite"
)

type InventoryServiceTest struct {
	suite.Suite

	repo    *mocks.PartRepo
	service PartService
}

func (s *InventoryServiceTest) SetupTest() {
	s.repo = mocks.NewPartRepo(s.T())
	s.service = NewPartService(s.repo)
}

func (s *InventoryServiceTest) TestGet() {
	ctx := context.Background()
	uuid := "engine-1"

	expected := &inventorypb.Part{Uuid: uuid}

	s.repo.
		On("Get", ctx, uuid).
		Return(expected, nil)

	res, err := s.service.Get(ctx, uuid)

	s.NoError(err)
	s.Equal(expected, res)

	s.repo.AssertExpectations(s.T())
}

func (s *InventoryServiceTest) TestGet_NotFound() {
	ctx := context.Background()
	uuid := "engine-1"

	s.repo.On("Get", ctx, uuid).Return(nil, errors.New("not found"))
	_, err := s.service.Get(ctx, uuid)

	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

func (s *InventoryServiceTest) TestList() {
	ctx := context.Background()
	filter := &inventorypb.PartsFilter{
		Uuids: []string{"engine-1", "wing-1"},
	}
	expected := []*inventorypb.Part{
		{Uuid: "engine-1", Name: "Engine"},
		{Uuid: "wing-1", Name: "Wing"},
	}
	s.repo.On("List", ctx, filter).Return(expected, nil)
	res, err := s.service.List(ctx, filter)
	s.NoError(err)
	s.Equal(res, expected)
	s.repo.AssertExpectations(s.T())

}

func TestInventoryServiceTest(t *testing.T) {
	suite.Run(t, new(InventoryServiceTest))
}

