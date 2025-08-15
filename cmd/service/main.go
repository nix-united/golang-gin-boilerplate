package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nix-united/golang-gin-boilerplate/docs"
	"github.com/nix-united/golang-gin-boilerplate/internal/config"
	"github.com/nix-united/golang-gin-boilerplate/internal/db"
	"github.com/nix-united/golang-gin-boilerplate/internal/handler"
	"github.com/nix-united/golang-gin-boilerplate/internal/provider"
	"github.com/nix-united/golang-gin-boilerplate/internal/repository"
	"github.com/nix-united/golang-gin-boilerplate/internal/server"
	"github.com/nix-united/golang-gin-boilerplate/internal/service/post"
	"github.com/nix-united/golang-gin-boilerplate/internal/service/user"
	"github.com/nix-united/golang-gin-boilerplate/internal/utils"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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

	var cfg config.ApplicationConfig
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port)

	// DB initialization
	gormDB, sqlDB, err := db.NewDBConnection(cfg.DB)
	if err != nil {
		return fmt.Errorf("new db connection: %w", err)
	}

	// Repository initialization
	userRepo := repository.NewUserRepository(gormDB)
	postRepo := repository.NewPostRepository(gormDB)

	// Services initialization
	userService := user.NewService(userRepo, utils.NewBcryptEncoder(bcrypt.DefaultCost))
	postService := post.NewService(postRepo)

	// Handlers initialization
	homeHandler := handler.NewHomeHandler()
	postHandler := handler.NewPostHandler(postService)
	authHandler := handler.NewAuthHandler(userService)
	jwtAuth := provider.NewJwtAuth(gormDB)

	// HTTP Server initialization
	httpServer := server.NewServer(cfg.HTTPServer, server.Handlers{
		HomeHandler:       homeHandler,
		AuthHandler:       authHandler,
		PostHandler:       postHandler,
		JwtAuthMiddleware: jwtAuth,
	})
	go func() {
		if err := httpServer.Run(); err != nil {
			log.Fatal("Server error: " + err.Error())
		}
	}()

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	<-shutdownChannel

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ApplicationShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("close db connection: %w", err)
	}

	return nil
}
