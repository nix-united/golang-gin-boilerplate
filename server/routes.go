package server

import (
	"basic_server/server/handler"
	"basic_server/server/provider"
)

func ConfigureRoutes(server *Server) {
	homeHandler := handler.HomeHandler{}
	postHandler := handler.PostHandler{DB: server.db}
	registerHandler := handler.RegisterHandler{DB: server.db}

	jwtAuth := provider.NewJwtAuth(server.db)

	server.engine.POST("/users", registerHandler.Register())
	server.engine.POST("/login", jwtAuth.Middleware().LoginHandler)

	needsAuth := server.engine.Group("/").Use(jwtAuth.Middleware().MiddlewareFunc())

	{
		needsAuth.GET("/", homeHandler.Index())
		needsAuth.GET("/refresh", jwtAuth.Middleware().RefreshHandler)
		needsAuth.POST("/posts", postHandler.SavePost())
		needsAuth.GET("/posts", postHandler.GetPosts())
	}
}
