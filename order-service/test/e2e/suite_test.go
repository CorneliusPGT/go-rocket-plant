package e2e

import (
	"context"
	"fmt"
	"net/http/httptest"
	"order-service/internal/handlers"
	"order-service/internal/oapi"
	repository "order-service/internal/repository/order"
	"order-service/internal/service/order"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type OrderE2ESuite struct {
	suite.Suite
	Env    *TestEnv
	Pool   *pgxpool.Pool
	Server *httptest.Server
	Client *oapi.Client
}

func (s *OrderE2ESuite) SetupSuite() {
	s.Env = SetupTestEnv(s.T())
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		s.Env.DbUser,
		s.Env.DbPassword,
		s.Env.DbHost,
		s.Env.DbPort,
		s.Env.DbName)
	pool, err := pgxpool.New(context.Background(), dsn)
	s.Require().NoError(err)
	s.Pool = pool

	repo := repository.NewRepository(pool)
	svc := order.NewService(repo, s.Env.InvMock, s.Env.PayMock)
	handler := &handlers.OrderHandler{Service: svc}

	ogenServer, err := oapi.NewServer(handler)
	s.Require().NoError(err)

	httpServer := httptest.NewServer(ogenServer)
	s.Server = httpServer

	client, err := oapi.NewClient(httpServer.URL)
	s.Require().NoError(err)
	s.Client = client
}

func (s *OrderE2ESuite) TearDownSuite() {
	if s.Server != nil {
		s.Server.Close()
	}
	if s.Pool != nil {
		s.Pool.Close()
	}
	if s.Env != nil {
		s.Env.TeardownTestEnv()
	}
}

func (s *OrderE2ESuite) SetupTest() {
	_, err := s.Pool.Exec(context.Background(), "TRUNCATE orders, order_items RESTART IDENTITY CASCADE")
	s.Require().NoError(err)
	s.Env.InvMock.ExpectedCalls = nil
	s.Env.InvMock.Calls = nil
}

func TestOrderE2E(t *testing.T) {
	suite.Run(t, new(OrderE2ESuite))
}
