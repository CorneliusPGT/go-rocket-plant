package e2e

import (
	"context"
	"fmt"
	"order-service/internal/migrator"
	"order-service/internal/mocks"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestEnv struct {
	Container  testcontainers.Container
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	InvMock    *mocks.InventoryService
	PayMock    *mocks.PaymentService
}

func SetupTestEnv(t *testing.T) *TestEnv {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "order",
			"POSTGRES_PASSWORD": "order",
			"POSTGRES_DB":       "orders",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")
	env := &TestEnv{
		Container:  container,
		DbHost:     host,
		DbPort:     port.Port(),
		DbUser:     "order",
		DbPassword: "order",
		DbName:     "orders",

		InvMock: mocks.NewInventoryService(t),
		PayMock: mocks.NewPaymentService(t),
	}
	runMigrations(env, t)
	return env
}

func runMigrations(env *TestEnv, t *testing.T) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", env.DbUser, env.DbPassword, env.DbHost, env.DbPort, env.DbName)

	migrationDir := "../../migrations"
	if err := migrator.Run(dsn, migrationDir); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}
}

func (t *TestEnv) TeardownTestEnv() {
	ctx := context.Background()
	t.Container.Terminate(ctx)
}
