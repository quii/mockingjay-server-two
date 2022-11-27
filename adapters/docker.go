package adapters

import (
	"context"
	"fmt"
	"io"
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

func StartMockingjayServer(
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
		ProviderType:     0,
		Logger:           nil,
		Reuse:            false,
	})

	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})
}

func RunMockingjayCLI(endpointDir string) (string, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		FromDockerfile: newTCDockerfile(),
		Cmd:            []string{"./svr", "-cdc=true", "-endpoints=/tmp/testresources/"},
		Mounts:         testcontainers.Mounts(testcontainers.BindMount(endpointDir, "/tmp/testresources")),
		WaitingFor:     wait.ForExit(),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return "", err
	}

	closer, err := container.Logs(ctx)
	if err != nil {
		return "", fmt.Errorf("couldn't get logs %w", err)
	}
	all, err := io.ReadAll(closer)
	if err != nil {
		return "", fmt.Errorf("problem reading logs %w", err)
	}
	return string(all), nil
}

func newTCDockerfile() testcontainers.FromDockerfile {
	return testcontainers.FromDockerfile{
		Context:       "../../.",
		Dockerfile:    dockerfileName,
		PrintBuildLog: true,
	}
}
