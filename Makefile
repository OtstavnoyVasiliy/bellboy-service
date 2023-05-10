.PHONY: migrate-up

DATABASE_CONFIG_FILE := config.json

# Переменные для подключения к базе данных, которые мы будем получать из файла JSON
DB_HOST := $(shell cat $(DATABASE_CONFIG_FILE) | jq -r '.database.host')
DB_PORT := $(shell cat $(DATABASE_CONFIG_FILE) | jq -r '.database.port')
DB_USER := $(shell cat $(DATABASE_CONFIG_FILE) | jq -r '.database.username')
DB_PASSWORD := $(shell cat $(DATABASE_CONFIG_FILE) | jq -r '.database.password')
DB_NAME := $(shell cat $(DATABASE_CONFIG_FILE) | jq -r '.database.name')
DB_DRIVER := $(shell cat $(DATABASE_CONFIG_FILE) | jq -r '.database.driver')
MG_DIRECTORY := $(shell cat $(DATABASE_CONFIG_FILE) | jq -r '.database.migrationsDir')

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

build:
	docker-compose build && docker-compose up -d

commit:
	git add .
	git commit -m "[$(BRANCH)] $(m)"
	git push

migrate-up:
	goose -dir $(MG_DIRECTORY) $(DB_DRIVER) "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

migrate-down:
	goose -dir $(MG_DIRECTORY) $(DB_DRIVER) "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" down