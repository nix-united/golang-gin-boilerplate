package main

import (
	"fmt"
	"log/slog"

	"github.com/nix-united/golang-gin-boilerplate/docs"
	application "github.com/nix-united/golang-gin-boilerplate/internal"
	"github.com/nix-united/golang-gin-boilerplate/internal/config"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

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

	var c config.Config
	if err := env.Parse(&c); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", c.HTTP.Host, c.HTTP.Port)
	application.Start(c)

	return nil
}
