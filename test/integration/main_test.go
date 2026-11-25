package integration_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"slices"
	"testing"
	"time"

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

const setupTimeout = 15 * time.Minute

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	ctx, cancel := context.WithTimeout(context.Background(), setupTimeout)
	defer cancel()

	startedAt := time.Now()

	shutdown, err := setupMain(ctx)
	if err != nil {
		attrs := []any{"err", err, "elapsed", time.Since(startedAt).Round(time.Second).String()}
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			attrs = append(attrs, "hint", fmt.Sprintf("setup exceeded its %s budget; see the container logs above", setupTimeout))
		}

		slog.Error("Failed to setup integration tests", attrs...)

		return 1
	}

	code := m.Run()

	cleanupCtx, cancelCleanup := setup.CleanupContext(ctx)
	defer cancelCleanup()

	if err := shutdown(cleanupCtx); err != nil {
		slog.Error("Failed to shutdown integration tests", "err", err)

		if code == 0 {
			code = 1
		}
	}

	return code
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

		cleanupCtx, cancel := setup.CleanupContext(ctx)
		defer cancel()

		if errShutdown := shutdown(cleanupCtx); errShutdown != nil {
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
		Driver:   "mysql",
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
