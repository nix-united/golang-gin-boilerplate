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
	server.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Repository Initialization
	userRepo := repository.NewUserRepository(server.DB)
	postRepo := repository.NewPostRepository(server.DB)

	// Services initialization
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo)

	// Handlers initialization
	homeHandler := handler.NewHomeHandler()
	postHandler := handler.NewPostHandler(postService)
	authHandler := handler.NewAuthHandler(userService)

	// Routes initialization
	server.Gin.POST("/users", authHandler.RegisterUser)

	jwtAuth := provider.NewJwtAuth(server.DB)
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
