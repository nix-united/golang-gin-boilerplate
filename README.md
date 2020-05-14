# Basic Project Structure

## main.go
Entry point of whole project is main.go file. Here you should init DB connection, init Server  and start listen for new connections.

## server package
`server` package will contain implementation of your server(requests, responses, handlers, routes, models).


Check `README.md` in `server` package to get more details.

## Swagger documentation


### Installation

1 Get the binary for github.com/swaggo/swag/cmd/swag:


    go get github.com/swaggo/swag/cmd/swag


2 In .env file set the value for the "HOST" variable. This is a host to which Swagger will make API requests. For example, for local development:


    HOST=localhost 
  
    
3 Run "swag init" in the project's root folder which contains the main.go file. This will parse your comments and generate the required files (docs folder and docs/docs.go).


     $GOPATH/bin/swag init 

    
### Usage

1. Run your app, and browse to {HOST}:{PORT}/swagger/index.html. You will see Swagger 2.0 API documents.
2. Using the API documentation, make requests to register a user (if necessary) and login.
3. After the successful login, copy a token from the response, then click "Authorize" and in a popup that opened, enter the value for "apiKey" in a form:
"Bearer {token}". For example:


    Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk0NDA5NjYsIm9yaWdfaWF0IjoxNTg5NDM5OTY2LCJ1c2VyX2lkIjo1fQ.f8dSG3NxFLHwyA5-XIYALT5GtXm4eiH-motqtqAUBOI 

   
Then, click "Authorize" and close the popup.
Now, you are able to make requests which require authentication.
   