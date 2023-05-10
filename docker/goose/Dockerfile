FROM alpine:3.15.7

WORKDIR /goose

RUN apk update && \
    apk add --update --no-cache make && \
    wget -O goose https://github.com/pressly/goose/releases/download/v3.10.0/goose_linux_x86_64 && \
    chmod +x goose && \
    mv goose /usr/bin/goose && \
    wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 && \
    chmod +x jq && \
    mv jq /usr/bin/jq

ADD migrations migrations
ADD config.json config.json
ADD Makefile Makefile

CMD make migrate-up