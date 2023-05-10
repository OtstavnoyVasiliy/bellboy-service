##
## BUILD
##
FROM golang:1.19 as build

WORKDIR /app

COPY go.* ./

RUN go mod download

ADD cmd ./cmd
ADD pkg ./pkg

RUN go build -o ./server ./cmd/server/main.go

##
## DEPLOY
##
FROM debian:11.6-slim

RUN apt-get update && apt-get install -y \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/server /app/server

ADD config.json config.json

EXPOSE 80

CMD ["/app/server"]