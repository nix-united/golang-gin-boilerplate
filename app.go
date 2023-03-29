package application

import (
	"basic_server/config"
	"basic_server/server"
	"log"
)

func Start(cfg *config.Config) {
	app := server.NewServer(cfg)

	server.ConfigureRoutes(app)

	err := app.Run(cfg.HTTP.Port)
	if err != nil {
		log.Fatal("Port already used")
	}
}
