include tools.mk

.PHONY: all build test test-basics test-bank test-challenges lint fmt bench tidy db-up db-down migrate help

all: tidy lint build test

build: ## Build all binaries (hello, bank-api, bank-cli)
	go build ./cmd/hello/...
	go build ./cmd/bank-api/...
	go build ./cmd/bank-cli/...

test: ## Run all tests
	go test ./...

test-basics: ## Run module 2 (Go basics) tests
	go test ./internal/basics/...

test-bank: ## Run module 3 (Go Bank) tests
	go test ./internal/bank/...

test-challenges: ## Run all challenge tests
	go test ./internal/challenges/...

lint: $(GOLANGCI_LINT) ## Run linter
	$(GOLANGCI_LINT) run ./...

fmt: ## Format Go code
	gofmt -w .

bench: ## Run all benchmarks
	go test -bench=. -benchmem ./...

tidy: ## Tidy go.mod dependencies
	go mod tidy

db-up: ## Start PostgreSQL database
	docker compose up -d postgres

db-down: ## Stop PostgreSQL database
	docker compose down

migrate: ## Run SQL migrations (instructions only)
	@echo "Migration tool not yet configured. See migration/ directory for SQL files."
	@echo "Recommended: use golang-migrate/migrate or goose."

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
