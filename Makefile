.PHONY: all build test test-basics test-bank test-challenges lint fmt bench tidy db-up db-down migrate

GOBIN ?= $(shell go env GOPATH)/bin

all: tidy lint build test

build:
	go build ./cmd/hello/...
	go build ./cmd/bank-api/...
	go build ./cmd/bank-cli/...

test:
	go test ./...

test-basics:
	go test ./internal/basics/...

test-bank:
	go test ./internal/bank/...

test-challenges:
	go test ./internal/challenges/...

lint:
	$(GOBIN)/golangci-lint run ./...

fmt:
	gofmt -w .

bench:
	go test -bench=. -benchmem ./...

tidy:
	go mod tidy

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

migrate:
	@echo "Migration tool not yet configured. See migration/ directory for SQL files."
	@echo "Recommended: use golang-migrate/migrate or goose."
