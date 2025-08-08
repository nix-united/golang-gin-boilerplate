package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nix-united/golang-gin-boilerplate/docs"
	"github.com/nix-united/golang-gin-boilerplate/internal/config"
	"github.com/nix-united/golang-gin-boilerplate/internal/server"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const shutdownTimeout = 5 * time.Minute

// @title Gin Demo App
// @version 1.0
// @description This is a demo version of Gin app.

// @contact.name NIX Solutions
// @contact.url https://www.nixsolutions.com/
// @contact.email ask@nixsolutions.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /
func main() {
	if err := run(); err != nil {
		slog.Error("Service run error", "err", err.Error())
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("load env file: %w", err)
	}

	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	app := server.NewServer(cfg)

	server.ConfigureRoutes(app)

	go func() {
		if err := app.Run(); err != nil {
			log.Fatal("Server error: " + err.Error())
		}
	}()

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	<-shutdownChannel

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	dbConnection, err := app.DB.DB()
	if err != nil {
		return fmt.Errorf("get db connection: %w", err)
	}

	if err := dbConnection.Close(); err != nil {
		return fmt.Errorf("close db connection: %w", err)
	}

	return nil
}
