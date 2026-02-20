package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestEnv struct {
	Container testcontainers.Container
	URI       string
}

func SetupTestEnv(t *testing.T) *TestEnv {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:7",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start mongo container: %v", err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "27017")

	uri := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	return &TestEnv{
		Container: container,
		URI:       uri,
	}
}

func (s *TestEnv) Teardown() {
	ctx := context.Background()
	s.Container.Terminate(ctx)
}
