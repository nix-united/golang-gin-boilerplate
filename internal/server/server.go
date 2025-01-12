package server

import (
	"github.com/nix-united/golang-gin-boilerplate/internal/config"
	"github.com/nix-united/golang-gin-boilerplate/internal/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	Cfg config.Config
	Gin *gin.Engine
	DB  *gorm.DB
}

func NewServer(cfg config.Config) *Server {
	return &Server{
		Cfg: cfg,
		Gin: gin.Default(),
		DB:  db.InitDB(cfg.DB),
	}
}

func (server *Server) Run(addr string) error {
	return server.Gin.Run(":" + addr)
}
