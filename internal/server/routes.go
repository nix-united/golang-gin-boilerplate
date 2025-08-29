package server

import (
	"net/http"

	"github.com/nix-united/golang-gin-boilerplate/internal/provider"
	"github.com/nix-united/golang-gin-boilerplate/internal/server/handler"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	HomeHandler       *handler.HomeHandler
	AuthHandler       *handler.AuthHandler
	PostHandler       *handler.PostHandler
	JwtAuthMiddleware provider.JwtAuthMiddleware
}

func configureRoutes(handlers Handlers) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	// Technical API route initialization
	// These endpoints exist solely to keep the service running and must not include any
	// business or processing logic.
	engine := gin.Default()

	engine.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	engine.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	api := engine.Group("/")

	// Private API routes initialization
	// These endpoints are used primarily for authentication/authorization and may carry sensitive data.
	// Do NOT log request or response bodies; doing so could expose client information.
	privateAPI := engine.Group("/")

	privateAPI.POST("/users", handlers.AuthHandler.RegisterUser)
	privateAPI.POST("/login", handlers.JwtAuthMiddleware.Middleware().LoginHandler)
	privateAPI.GET(
		"/refresh",
		handlers.JwtAuthMiddleware.Middleware().MiddlewareFunc(),
		handlers.JwtAuthMiddleware.Middleware().RefreshHandler,
	)

	// Authorized API route initialization
	//
	// These endpoints implement the core application logic and require authentication
	// before they can be accessed.
	authorizedAPI := api.Group("/", handlers.JwtAuthMiddleware.Middleware().MiddlewareFunc())

	authorizedAPI.GET("/", handlers.HomeHandler.Index)
	authorizedAPI.POST("/posts", handlers.PostHandler.SavePost)
	authorizedAPI.GET("/posts", handlers.PostHandler.GetPosts)
	authorizedAPI.GET("/post/:id", handlers.PostHandler.GetPostByID)
	authorizedAPI.PUT("/post/:id", handlers.PostHandler.UpdatePost)
	authorizedAPI.DELETE("/post/:id", handlers.PostHandler.DeletePost)

	return engine
}
