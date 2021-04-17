package server

import (
	"basic_server/handler"
	"basic_server/provider"
	"basic_server/repository"
	"basic_server/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {
	homeHandler := handler.HomeHandler{}
	postHandler := handler.PostHandler{DB: server.DB}
	registerHandler := handler.NewRegisterHandler()

	jwtAuth := provider.NewJwtAuth(server.DB)

	server.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.Gin.POST(
		"/users",
		registerHandler.RegisterUser(service.NewUserService(repository.NewUserRepository(server.DB))),
	)

	server.Gin.POST("/login", jwtAuth.Middleware().LoginHandler)

	needsAuth := server.Gin.Group("/").Use(jwtAuth.Middleware().MiddlewareFunc())
	needsAuth.GET("/", homeHandler.Index())
	needsAuth.GET("/refresh", jwtAuth.Middleware().RefreshHandler)
	needsAuth.POST("/posts", postHandler.SavePost)
	needsAuth.GET("/posts", postHandler.GetPosts)
	needsAuth.GET("/post/:id", postHandler.GetPostByID)
	needsAuth.PUT("/post/:id", postHandler.UpdatePost)
	needsAuth.DELETE("/post/:id", postHandler.DeletePost)
}
