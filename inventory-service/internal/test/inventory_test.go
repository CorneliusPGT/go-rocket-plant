package test

import (
	"context"
	"inventory-service/grpc/inventorypb"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *InvE2ESuite) TestListParts_Success() {
	_, err := s.Col.InsertMany(context.Background(), []interface{}{
		bson.M{
			"uuid":           "engine-1",
			"name":           "Main Engine",
			"price":          100.00,
			"stock_quantity": 10,
			"category":       1,
		},
		bson.M{
			"uuid":           "wing-1",
			"name":           "Wing",
			"price":          150.00,
			"stock_quantity": 5,
			"category":       1,
		},
		bson.M{
			"uuid":           "wing-2",
			"name":           "Wing 2",
			"price":          150.00,
			"stock_quantity": 5,
			"category":       1,
		},
		bson.M{
			"uuid":           "engine-2",
			"name":           "engine 2",
			"price":          120.00,
			"stock_quantity": 3,
			"category":       1,
		},
	})
	s.Require().NoError(err)
	resp, err := s.Client.ListParts(context.Background(), &inventorypb.ListPartsRequest{
		Filter: &inventorypb.PartsFilter{
			Uuids: []string{"engine-1", "wing-1"},
		},
	})
	s.Require().NoError(err)
	uuids := []string{resp.Parts[0].Uuid, resp.Parts[1].Uuid}
	s.Contains(uuids, "engine-1")
	s.Contains(uuids, "wing-1")
}

func (s *InvE2ESuite) TestListParts_EmptyResult() {
	resp, err := s.Client.ListParts(context.Background(), &inventorypb.ListPartsRequest{
		Filter: &inventorypb.PartsFilter{
			Uuids: []string{"engine-1"},
		},
	})
	s.Require().NoError(err)
	s.Equal(0, len(resp.Parts))

}

func (s *InvE2ESuite) TestGetPart_NotFound() {

	_, err := s.Client.GetPart(context.Background(), &inventorypb.GetPartRequest{
		Uuid: "engine1",
	})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "должна быть gRPC ошибка")
	s.Equal(codes.NotFound, st.Code())

}

func (s *InvE2ESuite) TestGetPart_Success() {
	_, err := s.Col.InsertOne(context.Background(), bson.M{
		"uuid":           "engine-1",
		"name":           "Main Engine",
		"price":          100.00,
		"stock_quantity": 10,
		"category":       1,
	})
	s.Require().NoError(err)

	resp, err := s.Client.GetPart(context.Background(), &inventorypb.GetPartRequest{
		Uuid: "engine-1",
	})

	s.Require().NoError(err)
	s.Equal("engine-1", resp.Part.Uuid)
}
