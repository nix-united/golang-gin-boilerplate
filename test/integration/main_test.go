package integration_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"slices"
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/config"
	"github.com/nix-united/golang-gin-boilerplate/internal/db"
	"github.com/nix-united/golang-gin-boilerplate/internal/slogx"
	"github.com/nix-united/golang-gin-boilerplate/test/setup"

	"gorm.io/gorm"
)

var (
	gormDB         *gorm.DB
	applicationURL *url.URL
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	shutdown, err := setupMain(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to setup integration tests", "err", err)
		return
	}

	m.Run()

	if err := shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown integration tests", "err", err)
		return
	}
}

func setupMain(ctx context.Context) (_ func(context.Context) error, err error) {
	err = slogx.Init(config.LoggerConfig{
		Application: "integration-tests",
		Level:       "DEBUG",
		AddSource:   true,
	})
	if err != nil {
		return nil, fmt.Errorf("init slog: %w", err)
	}

	shutdownCallbacks := make([]func(context.Context) error, 0)

	shutdown := func(ctx context.Context) error {
		var err error
		for _, callback := range slices.Backward(shutdownCallbacks) {
			err = errors.Join(err, callback(ctx))
		}

		if err != nil {
			return fmt.Errorf("shutdown callbacks: %w", err)
		}

		return nil
	}

	defer func() {
		if err == nil {
			return
		}

		if errShutdown := shutdown(context.WithoutCancel(ctx)); errShutdown != nil {
			err = errors.Join(err, errShutdown)
		}
	}()

	network, networkShutdown, err := setup.SetupNetwork(ctx)
	if err != nil {
		return nil, fmt.Errorf("setup network: %w", err)
	}

	shutdownCallbacks = append(shutdownCallbacks, networkShutdown)

	mysqlConfig, mysqlShutdown, err := setup.SetupMySQL(ctx, []string{network})
	if err != nil {
		return nil, fmt.Errorf("setup mysql: %w", err)
	}

	shutdownCallbacks = append(shutdownCallbacks, mysqlShutdown)

	applicationConfig, shutdownApplication, err := setup.SetupApplication(ctx, []string{network}, mysqlConfig)
	if err != nil {
		return nil, fmt.Errorf("setup application: %w", err)
	}

	shutdownCallbacks = append(shutdownCallbacks, shutdownApplication)

	applicationURL = applicationConfig.URL

	gdb, sqlDB, err := db.NewDBConnection(config.DBConfig{
		User:     mysqlConfig.User,
		Password: mysqlConfig.Password,
		Name:     mysqlConfig.Name,
		Host:     mysqlConfig.Host,
		Port:     mysqlConfig.ExposedPort,
	})
	if err != nil {
		return nil, fmt.Errorf("new gorm db connection: %w", err)
	}

	gormDB = gdb

	if err := db.Migrate(sqlDB); err != nil {
		return nil, fmt.Errorf("run db migrations: %w", err)
	}

	shutdownCallbacks = append(shutdownCallbacks, func(ctx context.Context) error {
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("close db connection: %w", err)
		}

		return nil
	})

	return shutdown, nil
}
