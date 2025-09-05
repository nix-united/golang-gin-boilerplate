package setup

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const appHTTPPort = "80"

type AppConfig struct {
	Port string
	Host string
	URL  *url.URL
}

func SetupApplication(
	ctx context.Context,
	networks []string,
	mySQLConfig MySQLConfig,
) (_ AppConfig, _ func(ctx context.Context) error, err error) {
	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				FromDockerfile: testcontainers.FromDockerfile{
					Context:    "../../",
					Dockerfile: "Dockerfile",
				},
				Env: map[string]string{
					"LOG_APPLICATION":     "golang-gin-boilerplate-integration-tests",
					"PORT":                appHTTPPort,
					"DB_DRIVER":           "mysql",
					"DB_USER":             mySQLConfig.User,
					"DB_PASSWORD":         mySQLConfig.Password,
					"DB_HOST":             mySQLConfig.ContainerName,
					"DB_PORT":             mySQLConfig.LocalPort,
					"DB_NAME":             mySQLConfig.Name,
					"JWT_SECRET":          "jwt-secret",
					"JWT_REALM":           "jwt-realm",
					"JWT_EXPIRATION_TIME": "300",
					"JWT_REFRESH_TIME":    "300",
				},
				WaitingFor: wait.
					ForAll(wait.ForHTTP("/health")).
					WithDeadline(time.Minute),
				Networks:     networks,
				ExposedPorts: []string{appHTTPPort},
			},
			Started: true,
		},
	)
	if err != nil {
		return AppConfig{}, nil, fmt.Errorf("generic container from app: %w", err)

	}

	shutdown := func(ctx context.Context) error {
		containerLogs, err := container.Logs(ctx)
		if err != nil {
			return fmt.Errorf("get application container logs: %w", err)
		}

		rawContainerLogs, err := io.ReadAll(containerLogs)
		if err != nil {
			return fmt.Errorf("read container logs: %w", err)
		}

		fmt.Println("\n\n\n### Start of application logs\n")
		fmt.Println(string(rawContainerLogs))
		fmt.Println("### End of application logs\n\n\n")

		if err := container.Terminate(ctx); err != nil {
			return fmt.Errorf("terminate app container: %w", err)
		}

		return nil
	}

	// This defer function exists to shutdown the container if any error occurs in next steps.
	defer func() {
		if err == nil {
			return
		}

		if errShutdown := shutdown(ctx); errShutdown != nil {
			err = errors.Join(err, errShutdown)
		}
	}()

	host, err := container.Host(ctx)
	if err != nil {
		return AppConfig{}, nil, fmt.Errorf("get host: %w", err)
	}

	port, err := container.MappedPort(ctx, appHTTPPort)
	if err != nil {
		return AppConfig{}, nil, fmt.Errorf("get exposed port: %w", err)
	}

	config := AppConfig{
		Port: port.Port(),
		Host: host,
		URL: &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%s", host, port.Port()),
		},
	}

	return config, shutdown, nil
}
