package server

import "basic_server/server/handler"

func ConfigureRoutes(server *Server)  {
	homeHandler := handler.HomeHandler{}
	postHandler := handler.PostHandler{DB: server.db}
	registerHandler := handler.RegisterHandler{DB: server.db}

	server.engine.POST("/users", registerHandler.Register())
	server.engine.GET("/", homeHandler.Index())
	server.engine.POST("/posts", postHandler.SavePost())
	server.engine.GET("/posts", postHandler.GetPosts())
}