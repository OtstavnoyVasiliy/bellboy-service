##
## BUILD
##
FROM golang:1.19 as build

WORKDIR /app

COPY go.* ./

RUN go mod download

ADD cmd ./cmd
ADD pkg ./pkg
ADD html ./html

RUN go build -o ./bot ./cmd/bot/main.go

##
## DEPLOY
##
FROM debian:11.6-slim

RUN apt-get update && apt-get install -y \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/bot /app/bot

ADD config.json config.json

CMD ["/app/bot"]