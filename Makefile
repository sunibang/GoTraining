include tools.mk

.PHONY: all build build-hello clean test test-hello test-basics test-bank test-challenges lint fmt bench tidy db-up db-down migrate docker-build-hello docker-run-hello help generate

HELLO_IMAGE ?= hello:latest

all: tools generate clean tidy lint build test

generate: $(MOCKERY) ## Generate mocks
	$(MOCKERY)

clean: ## Remove built binaries
	rm -rf bin/
	go clean -testcache

build: generate ## Build all binaries (hello, bank-api, bank-cli)
	@mkdir -p bin
	go build -o bin/hello ./cmd/hello/main.go
	go build -o bin/bank-api ./cmd/bank-api/main.go
	go build -o bin/bank-cli ./cmd/bank-cli/main.go
	@chmod +x bin/*

build-hello: ## Build hello world binaries
	# Building production ready executable
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o ./bin/hello ./cmd/hello/main.go

test: generate ## Run all tests
	go test `go list ./... | grep -v 'api/transfer'`

test-hello:  ## Run all hello world tests
	go test ./cmd/hello/...
	go test ./internal/hello/...

test-basics: ## Run module 2 (Go basics) tests
	go test ./internal/basics/...

test-bank: generate ## Run module 3 (Go Bank) tests
	go test `go list ./internal/bank/... | grep -v 'api/transfer'`

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

docker-build-hello: ## Build docker image for hello world
	docker build -f ./cmd/hello/Dockerfile -t $(HELLO_IMAGE) .

docker-run-hello: ## Run hello world through docker
	docker run --rm $(HELLO_IMAGE) $(NAME)

db-up: ## Start PostgreSQL database
	docker-compose up -d postgres

db-down: ## Stop PostgreSQL database
	docker-compose down

migrate: ## Run SQL migrations (instructions only)
	@echo "Migration tool not yet configured. See migration/ directory for SQL files."
	@echo "Recommended: use golang-migrate/migrate or goose."

# Thin-slice targets
run-bank-api: build ## Start bank API server
	./bin/bank-api

run-bank-cli: build ## Use bank CLI (example: make run-bank-cli ARGS="account create 'John Doe'")
	./bin/bank-cli $(ARGS)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
