# Temporal targets

.PHONY: temporal-up temporal-down worker-start test-temporal generate.mocks

temporal-up: ## Start Temporal dev server and WireMock
	docker compose up -d temporal wiremock

temporal-down: ## Stop Temporal and WireMock
	docker compose down temporal wiremock

worker-start: ## Start the Temporal order processing worker
	go run ./cmd/worker/... -config="./config/worker/local/config.yaml"

test-temporal: ## Run module 4 (Temporal) tests
	go test ./internal/temporal/...

generate.mocks: $(MOCKGEN) ## Generate mocks
	$(MOCKGEN) -destination=internal/temporal/activities/mocks/mock_inventory_checker.go \
	           -package=mocks \
	           github.com/romangurevitch/go-training/internal/temporal/activities InventoryChecker
