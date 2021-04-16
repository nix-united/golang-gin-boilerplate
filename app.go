package application

import (
	"basic_server/config"
	"basic_server/server"
	"log"
)

func Start(cfg *config.Config) {
	app := server.NewServer(cfg)

	routes.ConfigureRoutes(app)

	err := app.Start(cfg.HTTP.Port)
	if err != nil {
		log.Fatal("Port already used")
	}
}
