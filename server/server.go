package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Server struct {
	Gin *gin.Engine
	DB  *gorm.DB
}

func NewServer(dbConnection *gorm.DB) *Server {
	return &Server{
		Gin: gin.Default(),
		DB:  dbConnection,
	}
}

func (server *Server) Run(addr string) error {
	return server.Gin.Run(":" + addr)
}

func (server *Server) Engine() *gin.Engine {
	return server.Gin
}

func (server *Server) Database() *gorm.DB {
	return server.DB
}
