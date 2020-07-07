package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Server struct {
	engine *gin.Engine
	db     *gorm.DB
}

func NewServer(dbConnection *gorm.DB) *Server {
	return &Server{
		engine: gin.Default(),
		db:     dbConnection,
	}
}

func (server *Server) Run(addr string) error {
	return server.engine.Run(":" + addr)
}

func (server *Server) Engine() *gin.Engine {
	return server.engine
}

func (server *Server) Database() *gorm.DB {
	return server.db
}
