package main

import (
	"basic_server/server"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := server.NewServer()
	server.ConfigureRoutes(app)
	app.Run(os.Getenv("PORT"))
}