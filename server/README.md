# Server package

## server.go
Main file here is `server.go`. Server struct defined there. Server struct stores deps for other parts of the app (handlers, services, etc).

Two methods in `server.go`: `Run()` and `NewServer()`.

`NewServer` - constructor.

`Run` - start accepting connections.


## routes.go
`routes.go` is a file where connection between handlers and routes are created.