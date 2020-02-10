package server

import (
	"basic_server/server/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Server struct {
	engine *gin.Engine
	db     *gorm.DB
}

func NewServer() *Server {
	return &Server{
		engine: gin.Default(),
		db:     db.InitDB(),
	}
}

func (server *Server) Run(addr string) error {
	return server.engine.Run(":" + addr)
}
