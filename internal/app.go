package application

import (
	"log"

	"github.com/nix-united/golang-gin-boilerplate/internal/config"
	"github.com/nix-united/golang-gin-boilerplate/internal/server"
)

func Start(cfg config.Config) {
	app := server.NewServer(cfg)

	server.ConfigureRoutes(app)

	err := app.Run(cfg.HTTP.Port)
	if err != nil {
		log.Fatal("Port already used")
	}
}
