package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nix-united/golang-gin-boilerplate/internal/config"
)

type Server struct {
	server *http.Server
}

func NewServer(cfg config.HTTPServerConfig, handlers Handlers) *Server {
	return &Server{
		server: &http.Server{
			Addr:              cfg.Host + ":" + cfg.Port,
			Handler:           configureRoutes(handlers),
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("run http server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	return nil
}
