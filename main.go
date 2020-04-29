package main

import (
	"log"
	"os"

	"basic_server/server"
	"basic_server/server/db"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connection := db.InitDB()

	defer func() {
		if err := connection.DB().Close(); err != nil {
			log.Fatal(err)
		}
	}()

	app := server.NewServer(connection)
	server.ConfigureRoutes(app)

	if err := app.Run(os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}
