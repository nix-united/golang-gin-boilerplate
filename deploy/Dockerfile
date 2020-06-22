# base image
FROM golang:1.13.3-alpine3.10 as build

# install system utils
RUN apk add --no-cache --update git build-base openssh-client

# change workdir
WORKDIR /go/src/api

## read build arguments
ARG WEB_PRIVATE_KEY
ARG GIT_DOMAIN
#
## configuring ssh client
RUN mkdir ~/.ssh && \
    echo "$WEB_PRIVATE_KEY" | tr -d '\r' > ~/.ssh/id_rsa && \
    chmod 600 ~/.ssh/id_rsa && \
    ssh-keyscan -H $GIT_DOMAIN >> ~/.ssh/known_hosts

# copy project
COPY . .

RUN  git config --global http.sslVerify true &&\
     # install swagger and initialize documentation
     go get -u github.com/pressly/goose/cmd/goose &&\
     go get -v github.com/swaggo/swag/cmd/swag &&\
     swag init &&\
     # install project dependecies
     go get -v ./... &&\
     # build application
     go build -o demo-gin

FROM alpine

WORKDIR /app

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

ENV DB_USERNAME ''
ENV DB_PASSWORD ''
ENV DB_NAME     ''
ENV DB_HOST     ''
ENV DB_PORT     ''

COPY --from=build /go/src/api/demo-gin /app/demo-gin
COPY --from=build /go/src/api/server/db/migrations /app/migrations
COPY --from=build /go/bin/goose /app/migrations

CMD /wait &&\
    ./migrations/goose -dir "./migrations" mysql "$DB_USERNAME:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true" up &&\
    ./demo-gin
