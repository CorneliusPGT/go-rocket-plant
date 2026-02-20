package test

import (
	"context"
	"inventory-service/grpc/handlers"
	"inventory-service/grpc/inventorypb"
	"inventory-service/internal/service"
	repo "inventory-service/repository"
	"net"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InvE2ESuite struct {
	suite.Suite
	Env      *TestEnv
	Mongo    *mongo.Client
	Col      *mongo.Collection
	Server   *grpc.Server
	Listener net.Listener
	Client   inventorypb.InventoryServiceClient
}

func (s *InvE2ESuite) SetupSuite() {
	ctx := context.Background()

	s.Env = SetupTestEnv(s.T())

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(s.Env.URI))
	s.Require().NoError(err)
	s.Mongo = client

	s.Col = client.Database("inventory_test").Collection("items")

	repo := repo.NewMongoRepo(s.Col)
	svc := service.NewPartService(repo)
	handler := handlers.NewInventoryHandler(svc)
	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)
	s.Listener = lis

	grpcServer := grpc.NewServer()
	inventorypb.RegisterInventoryServiceServer(grpcServer, handler)
	s.Server = grpcServer
	go grpcServer.Serve(lis)

	conn, err := grpc.NewClient(s.Listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.Client = inventorypb.NewInventoryServiceClient(conn)
}

func (s *InvE2ESuite) TearDownSuite() {
	if s.Server != nil {
		s.Server.GracefulStop()
	}
	if s.Listener != nil {
		s.Listener.Close()
	}
	if s.Mongo != nil {
		s.Mongo.Disconnect(context.Background())
	}
	if s.Env != nil {
		s.Env.Teardown()
	}
}

func (s *InvE2ESuite) SetupTest() {
	s.Col.Drop(context.Background())
}

func TestInventoryE2E(t *testing.T) {
	suite.Run(t, new(InvE2ESuite))
}
