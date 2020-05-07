package service

import (
	"basic_server/server"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var ts *testServer
var once sync.Once

type testServer struct {
	engine   *gin.Engine
	database *gorm.DB
}

func TestServer() *testServer {
	once.Do(func() {
		err := godotenv.Load("../.env")

		if err != nil {
			log.Fatal("Error loading .env file")
		}

		srv := server.NewServer()

		server.ConfigureRoutes(srv)

		ts = &testServer{
			engine:   srv.Engine(),
			database: srv.Database(),
		}
	})

	return ts
}

func (ts *testServer) Engine() *gin.Engine {
	return ts.engine
}

func (ts *testServer) DatabaseDriver() *gorm.DB {
	return ts.database
}
