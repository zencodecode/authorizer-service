ifneq (,$(wildcard .env))
	include .env
	export
endif

MIGRATE := $(HOME)/go/bin/migrate
VERSION := $(shell git describe --tags --abbrev=0 || echo "0.0.0")
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
DB_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=$(POSTGRES_SSLMODE)
MIGRATION_PATH := internal/infrastructure/persistence/postgres/migrations

local-http-private:
	air

migrate-create:
	$(MIGRATE) create -ext sql -dir $(MIGRATION_PATH) -seq $(MIGRATION_NAME)

migrate-up:
	$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DB_URL)" down

migrate-force:
	$(MIGRATE) -path $(MIGRATION_PATH) -database "$(DB_URL)" force $(ver)

build:
	go mod tidy
	go build -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)" -o ./bin/main ./cmd/api

run: build
	./bin/main

clean:
	rm -f ./bin/main