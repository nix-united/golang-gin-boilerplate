# Overview
It's an API Skeleton project based on Gin framework.
Our aim is reducing development time on default features that you can meet very often when your work on API.
There is a useful set of tools that described below. Feel free to contribute!

## What's inside:

- Registration
- Authentication with JWT
- CRUD API for posts
- Migrations
- Request validation
- Swagger docs
- Environment configuration
- Docker development environment

## Usage
1. Copy .env.dist to .env and set the environment variables.
2. Run your application using the command in the terminal:

    `docker-compose up`
3. Browse to {HOST}:{EXPOSE_PORT}/swagger/index.html. You will see Swagger 2.0 API documents.
4. Using the API documentation, make requests to register a user (if necessary) and login.
5. After the successful login, copy a token from the response, then click "Authorize" and in a popup that opened, enter the value for "apiKey" in a form:
"Bearer {token}". For example:


    Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk0NDA5NjYsIm9yaWdfaWF0IjoxNTg5NDM5OTY2LCJ1c2VyX2lkIjo1fQ.f8dSG3NxFLHwyA5-XIYALT5GtXm4eiH-motqtqAUBOI 

   
Then, click "Authorize" and close the popup.
Now, you are able to make requests which require authentication.

# Basic Project Structure

## main.go
Entry point of whole project is main.go file. Here you should init DB connection, init Server  and start listen for new connections.

## server package
`server` package will contain implementation of your server(requests, responses, handlers, routes, models).


Check `README.md` in `server` package to get more details.

## Code quality
For control code quality we are use [golangci-lint](https://github.com/golangci/golangci-lint).
Golangci-lint is a linters aggregator.

Why we use linters? Linters help us:
1. Finding critical bugs
2. Finding bugs before they go live
3. Finding performance errors
4. To speed up the code review, because reviewers do not spend time searching for syntax errors and searching for
violations of generally accepted code style
5. The quality of the code is guaranteed at a fairly high level.

### How to use
Linter tool wrapped to docker-compose and first of all need to build container with linters

- `make lint-build`

Next you need to run linter to check bugs ant errors

- `make lint-check` - it will log to console what bugs and errors linters found

Finally, you need to fix all problems manually or using autofixing (if it's supported by the linter)

- `make lint-fix` 
