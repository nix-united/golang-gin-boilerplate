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
	"github.com/nix-united/golang-gin-boilerplate/internal/provider"
	"github.com/nix-united/golang-gin-boilerplate/internal/repository"
	"github.com/nix-united/golang-gin-boilerplate/internal/server"
	"github.com/nix-united/golang-gin-boilerplate/internal/server/handler"
	"github.com/nix-united/golang-gin-boilerplate/internal/server/middleware"
	"github.com/nix-united/golang-gin-boilerplate/internal/service/post"
	"github.com/nix-united/golang-gin-boilerplate/internal/service/user"
	"github.com/nix-united/golang-gin-boilerplate/internal/slogx"
	"github.com/nix-united/golang-gin-boilerplate/internal/utils"

	"github.com/caarlos0/env"
	"github.com/google/uuid"
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

	if err := slogx.Init(cfg.Logger); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	traceStarter := slogx.NewTraceStarter(uuid.NewV7)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port)

	// DB initialization
	gormDB, sqlDB, err := db.NewDBConnection(cfg.DB)
	if err != nil {
		return fmt.Errorf("new db connection: %w", err)
	}

	// Repository initialization
	userRepository := repository.NewUserRepository(gormDB)
	postRepository := repository.NewPostRepository(gormDB)

	// Services initialization
	bcryptEncoder := utils.NewBcryptEncoder(bcrypt.DefaultCost)
	userService := user.NewService(userRepository, bcryptEncoder)
	postService := post.NewService(postRepository)

	// Handlers initialization
	homeHandler := handler.NewHomeHandler()
	postHandler := handler.NewPostHandler(postService)
	authHandler := handler.NewAuthHandler(userService)
	jwtAuth := provider.NewJwtAuth(gormDB)

	// Middlewares initialization
	requestLoggerMiddleware := middleware.NewRequestLoggerMiddleware(traceStarter)
	requestDebuggerMiddleware := middleware.NewRequestDebuggerMiddleware()

	// HTTP Server initialization
	routes := server.ConfigureRoutes(server.Handlers{
		HomeHandler:                homeHandler,
		AuthHandler:                authHandler,
		PostHandler:                postHandler,
		JwtAuthMiddleware:          jwtAuth,
		RequestLoggingMiddleware:   requestLoggerMiddleware.Handle,
		RequestDebuggingMiddleware: requestDebuggerMiddleware.Handle,
	})

	httpServer := server.NewServer(cfg.HTTPServer, routes)
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
