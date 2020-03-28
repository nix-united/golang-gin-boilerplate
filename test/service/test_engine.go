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
	databaseClener
	engine   *gin.Engine
	database *gorm.DB
}

func TestServer() *testServer {
	once.Do(func() {
		err := godotenv.Load("../.env.testing")

		if err != nil {
			log.Fatal("Error loading .env.testing file")
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
