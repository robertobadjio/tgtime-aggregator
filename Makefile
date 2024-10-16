#!/usr/bin/make
.DEFAULT_GOAL := help
.PHONY: help

DOCKER_COMPOSE ?= docker compose -f docker-compose.yml

include .env

export GOOS=linux
export GOARCH=amd64

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-deps: ## Install dependencies for protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps: ## Get dependencies for protobuf
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

fmt: ## Automatically format source code
	go fmt ./...
.PHONY:fmt

lint: fmt ## Check code (lint)
	golangci-lint run ./... --config .golangci.pipeline.yaml
.PHONY:lint

vet: fmt ## Check code (vet)
	go vet -vettool=$(which shadow) ./...
.PHONY:vet

vet-shadow: fmt ## Check code with detect shadow (vet)
	go vet -vettool=$(which shadow) ./...
.PHONY:vet

build: ## Build service containers
	$(DOCKER_COMPOSE) build

up: vet ## Start services
	$(DOCKER_COMPOSE) up -d $(SERVICES)

down: ## Down services
	$(DOCKER_COMPOSE) down

clean: ## Delete all containers
	$(DOCKER_COMPOSE) down --remove-orphans

migrate-local-up: ## Migrates the database to the latest version
	docker run -v "${DATABASE_MIGRATION_DIR}:/migrations" --network host migrate/migrate -path=/migrations/ -database postgres://${DATABASE_PG_USER}:${DATABASE_PG_PASSWORD}@127.0.0.1:${DATABASE_PG_PORT}/${DATABASE_PG_NAME}?sslmode=${DATABASE_PG_SSL_MODE} up

migrate-local-create: ## Creates a new migration file with the given name. Ex: name=create_users_table
	docker run -v "${DATABASE_MIGRATION_DIR}:/migrations" -ext sql --network host migrate/migrate -path=/migrations/ -database postgres://${DATABASE_PG_USER}:${DATABASE_PG_PASSWORD}@127.0.0.1:${DATABASE_PG_PORT}/${DATABASE_PG_NAME}?sslmode=${DATABASE_PG_SSL_MODE} create $(name)

migrate-local-down: ## Migrates the database down
	docker run -v "${DATABASE_MIGRATION_DIR}:/migrations" -ext sql --network host migrate/migrate -path=/migrations/ -database postgres://${DATABASE_PG_USER}:${DATABASE_PG_PASSWORD}@127.0.0.1:${DATABASE_PG_PORT}/${DATABASE_PG_NAME}?sslmode=${DATABASE_PG_SSL_MODE} down

generate-api:
	make generate-time-proto-api

generate-time-proto-api: ## Generate pb.go files for time API
	mkdir -p pkg/api/time_v1
	protoc --proto_path=api/v1/pb/time \
	--go_out=pkg/api/time_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/api/time_v1 --go-grpc_opt=paths=source_relative \
	time.proto