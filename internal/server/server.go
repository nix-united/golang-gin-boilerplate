package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nix-united/golang-gin-boilerplate/internal/config"
	"github.com/nix-united/golang-gin-boilerplate/internal/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	Cfg    config.Config
	Gin    *gin.Engine
	DB     *gorm.DB
	server *http.Server
}

func NewServer(cfg config.Config) *Server {
	engine := gin.Default()

	return &Server{
		Cfg: cfg,
		Gin: engine,
		DB:  db.InitDB(cfg.DB),
		server: &http.Server{
			Addr:              ":" + cfg.HTTP.Port,
			Handler:           engine,
			ReadHeaderTimeout: 10 * time.Minute,
			ReadTimeout:       10 * time.Minute,
			WriteTimeout:      10 * time.Minute,
		},
	}
}

func (server *Server) Run() error {
	if err := server.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run http server: %w", err)
	}

	return nil
}

func (server *Server) Shutdown(ctx context.Context) error {
	if err := server.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	return nil
}
