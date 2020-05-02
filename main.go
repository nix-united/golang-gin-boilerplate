package main

import (
	"basic_server/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := server.NewServer()
	server.ConfigureRoutes(app)
	err = app.Run(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
