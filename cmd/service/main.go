package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nix-united/golang-gin-boilerplate/docs"
	application "github.com/nix-united/golang-gin-boilerplate/internal"
	"github.com/nix-united/golang-gin-boilerplate/internal/config"

	"github.com/joho/godotenv"
)

// @title Gin Demo App
// @version 1.0
// @description This is a demo version of Gin app.

// @contact.name NIX Solutions
// @contact.url https://www.nixsolutions.com/
// @contact.email ask@nixsolutions.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	application.Start(config.NewConfig())
}
