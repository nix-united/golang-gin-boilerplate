# Server package

## server.go
Main file here is `server.go`. Server struct defined there. Server struct stores deps for other parts of the app (handlers, services, etc).

Two methods in `server.go`: `Run()` and `NewServer()`.

`NewServer` - constructor.

`Run` - start accepting connections.


## routes.go
`routes.go` is a file where connection between handlers and routes are created.

## db package
`db` package is responsible for connection with a database. Place for DB connection initiation.
 
 ## handler
 `handler` package is responsible for connection of business logic layer and transport layer.

## model
`model` package is a place where to store app models.

## repository
`repository` package is a place where to store repositories.

## request
`request` package is a place where to store requests.

## response
`response` package is a place where to store requests.

## service
`service` package is a place where to store business logic of the app.

It is also possible to add other packages here(`utils`, for example).