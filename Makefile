.PHONY: dc run test lint

Red='\033[0;31m'
Green='\033[0;32m'
Color_Off='\033[0m'

help:
	@echo ${Red}"Please select a subcommand:"${Color_Off}
	@echo ${Green}"make run-postgres"${Color_Off}" to run postgres"
	@echo ${Green}"make run"${Color_Off}" to run bookshop"
	@echo
	@echo ${Red}"Or use docker-compose:"
	@echo ${Green}"make dc"${Color_Off}" to run docker-compose"
	@echo ${Green}"make down"${Color_Off}" to stop docker-compose"
	@echo
	@echo ${Green}"make lint"${Color_Off}" to run linter"
	@echo ${Green}"make test"${Color_Off}" to run unit tests"

dc:
	@docker-compose -f ./docker-compose.yml up --remove-orphans --build

build:
	go build -race -o app cmd/main.go

down:
	@docker-compose -f ./docker-compose.yml down

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.59.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...



run-postgres:
	docker run --name pg -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=bookshop -e PGDATA=/var/lib/postgresql/data/pgdata -p 5433:5432 -v pg_data:/var/lib/postgresql/data -d postgres:latest

migrate-up:
	migrate -path ./internal/app/migrations -database "postgres://postgres:password@localhost:5433/bookshop?sslmode=disable" up

run:
	go build -race -o app cmd/main.go && \
	HTTP_ADDR=:8080 \
	DEBUG_ERRORS=1 \
	DSN="postgres://postgres:password@127.0.0.1:5433/bookshop?sslmode=disable" \
	MIGRATIONS_PATH="file://./internal/app/migrations" \
	./app

