package service

import (
	"basic_server/server"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ts *server.Server
var once sync.Once

func TestEngine() *gin.Engine {
	once.Do(func() {
		err := godotenv.Load("../.env.testing")

		if err != nil {
			log.Fatal("Error loading .env.testing file")
		}

		ts = server.NewServer()

		server.ConfigureRoutes(ts)
	})

	return ts.Engine()
}
