package testhelpers

import (
	"context"
	"path/filepath"
	"runtime"
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo"
	"github.com/docker/go-connections/nat"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const containerName = "centrifugo_test"

func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	// From backend/tests/testhelpers/centrifugo.go -> project root
	return filepath.Join(filepath.Dir(filename), "..", "..", "..")
}

func SetupCentrifugoClient(ctx context.Context) (client *centrifugo.CentrifugoClient, cleanup func(), err error) {
	projectRoot := getProjectRoot()

	envMap, err := godotenv.Read(filepath.Join(projectRoot, ".env"))
	if err != nil {
		return nil, nil, err
	}

	configPath := filepath.Join(projectRoot, "centrifugo", "config.json")

	grpcPort := envMap["CENTRI_GRPC_API_PORT"]
	natPort := nat.Port(grpcPort + "/tcp")

	container, err := testcontainers.Run(ctx, "centrifugo/centrifugo:v6",
		testcontainers.WithCmd("centrifugo", "-c", "config.json"),
		testcontainers.WithFiles(testcontainers.ContainerFile{
			HostFilePath:      configPath,
			ContainerFilePath: "/config.json",
		}),
		testcontainers.WithName(containerName),
		testcontainers.WithEnv(envMap),
		testcontainers.WithExposedPorts(grpcPort+"/tcp"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort(natPort).WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return nil, nil, err
	}

	// Get the mapped host port
	mappedPort, err := container.MappedPort(ctx, natPort)
	if err != nil {
		return nil, nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	client, err = centrifugo.NewCentrifugoClient(host, int(mappedPort.Int()))
	if err != nil {
		return nil, nil, err
	}

	return client, func() {
		container.Terminate(ctx)
	}, nil
}

type CentrifugoConfig struct {
	Host string
	Port int
}

func StartCentrifugoContainer(ctx context.Context) (*CentrifugoConfig, func(), error) {
	projectRoot := getProjectRoot()

	// 1. Load Env
	envMap, err := godotenv.Read(filepath.Join(projectRoot, ".env"))
	if err != nil {
		return nil, nil, err
	}

	// 2. Identify Internal Port
	// We need to know which port inside the container is the gRPC port
	internalPortStr := envMap["CENTRI_GRPC_API_PORT"]

	// 3. Prepare Config File
	configPath := filepath.Join(projectRoot, "centrifugo", "config.json")

	// 4. Start Container
	req := testcontainers.ContainerRequest{
		Image: "centrifugo/centrifugo:v6",
		Cmd:   []string{"centrifugo", "-c", "/config.json"},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      configPath,
				ContainerFilePath: "/config.json",
			},
		},
		Env: envMap,
		// Vital: Expose the port so we can map it to localhost
		ExposedPorts: []string{internalPortStr + "/tcp"},
		// Vital: Wait for Centrifugo to actually accept connections
		WaitingFor: wait.ForListeningPort(nat.Port(internalPortStr + "/tcp")).
			WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	// 5. Cleanup Closure
	cleanup := func() {
		container.Terminate(ctx)
	}

	// 6. Get the Mapped Port (The "Real" Port on your laptop)
	// Since we run tests on Host, we connect to localhost, not the container name.
	mappedPort, err := container.MappedPort(ctx, nat.Port(internalPortStr+"/tcp"))
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	config := &CentrifugoConfig{
		Host: "localhost", // Access via localhost when running go test
		Port: mappedPort.Int(),
	}

	return config, cleanup, nil
}
