package server

import (
	"basic_server/server/handler"
	"basic_server/server/provider"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {
	homeHandler := handler.HomeHandler{}
	postHandler := handler.PostHandler{DB: server.db}
	registerHandler := handler.RegisterHandler{DB: server.db}

	jwtAuth := provider.NewJwtAuth(server.db)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.engine.POST("/users", registerHandler.Register())
	server.engine.POST("/login", jwtAuth.Middleware().LoginHandler)

	needsAuth := server.engine.Group("/").Use(jwtAuth.Middleware().MiddlewareFunc())

	{
		needsAuth.GET("/", homeHandler.Index())
		needsAuth.GET("/refresh", jwtAuth.Refresh)
		needsAuth.POST("/posts", postHandler.SavePost())
		needsAuth.GET("/posts", postHandler.GetPosts())
	}
}
