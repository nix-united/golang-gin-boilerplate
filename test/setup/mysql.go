package setup

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	mysqlImage         = "mysql:8.4.5"
	mysqlDatabase      = "db_name"
	mysqlUsername      = "username"
	mysqlPassword      = "password"
	mysqlPort          = "3306"
	mysqlHost          = "localhost"
	mysqlContainerName = "golang_gin_boilerplate_mysql_db"
)

type MySQLConfig struct {
	User          string
	Password      string
	Host          string
	ExposedPort   string
	LocalPort     string
	Name          string
	ContainerName string
}

func SetupMySQL(ctx context.Context, networks []string) (_ MySQLConfig, _ func(ctx context.Context) error, err error) {
	container, err := mysql.Run(
		ctx,
		mysqlImage,
		mysql.WithDatabase(mysqlDatabase),
		mysql.WithUsername(mysqlUsername),
		mysql.WithPassword(mysqlPassword),
		testcontainers.WithWaitStrategyAndDeadline(
			time.Minute,
			wait.ForLog("X Plugin ready for connections. Bind-address: '::' port: 33060, socket: /var/run/mysqld/mysqlx.sock"),
		),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name:     mysqlContainerName,
				Networks: networks,
			},
		}),
	)
	if err != nil {
		return MySQLConfig{}, nil, fmt.Errorf("run mysql container: %w", err)
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

		fmt.Print("\n\n\n### Start of database logs\n\n")
		fmt.Println(string(rawContainerLogs))
		fmt.Print("### End of database logs\n\n\n\n")

		if err := container.Terminate(ctx); err != nil {
			return fmt.Errorf("terminate mysql container: %w", err)
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

	port, err := container.MappedPort(ctx, mysqlPort+"/tcp")
	if err != nil {
		return MySQLConfig{}, nil, fmt.Errorf("get mysql exposed port: %w", err)
	}

	config := MySQLConfig{
		User:          mysqlUsername,
		Password:      mysqlPassword,
		Host:          mysqlHost,
		ExposedPort:   port.Port(),
		LocalPort:     mysqlPort,
		Name:          mysqlDatabase,
		ContainerName: mysqlContainerName,
	}

	return config, shutdown, nil
}
