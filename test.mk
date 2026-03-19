# Test targets

.PHONY: test test-hello test-basics test-bank test-challenges bench

test: generate ## Run all tests
	go test ./...

test-hello: ## Run all hello world tests
	go test ./cmd/hello/... ./internal/hello/...

test-basics: ## Run module 2 (Go basics) tests
	go test ./internal/basics/...

test-bank: generate ## Run module 3 (Go Bank) tests
	go test ./internal/bank/...

test-challenges: ## Run all challenge tests
	go test ./internal/challenges/...

bench: ## Run all benchmarks
	go test -bench=. -benchmem ./...
