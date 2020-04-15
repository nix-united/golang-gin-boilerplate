package main

import (
	"basic_server/server"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"basic_server/docs"
)

// @title Gin Demo App
// @version 1.0
// @description This is a demo version of Gin app.

// @contact.name NIX Solutions
// @contact.url https://www.nixsolutions.com/
// @contact.email ask@nixsolutions.com

// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	app := server.NewServer()
	server.ConfigureRoutes(app)
	app.Run(os.Getenv("PORT"))
}