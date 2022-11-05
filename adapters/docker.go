package adapters

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	startupTimeout = 5 * time.Second
	dockerfileName = "Dockerfile"
)

func StartDockerServer(
	t testing.TB,
	stubServerPort string,
	configServerPort string,
) {
	t.Helper()

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		FromDockerfile: newTCDockerfile(),
		ExposedPorts: []string{
			fmt.Sprintf("%s:%s", stubServerPort, stubServerPort),
			fmt.Sprintf("%s:%s", configServerPort, configServerPort),
		},
		WaitingFor: wait.ForListeningPort(nat.Port(stubServerPort)).WithStartupTimeout(startupTimeout),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})
}

func newTCDockerfile() testcontainers.FromDockerfile {
	return testcontainers.FromDockerfile{
		Context:       "../../.",
		Dockerfile:    dockerfileName,
		PrintBuildLog: true,
	}
}
