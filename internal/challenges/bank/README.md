# Go Bank Transfer Quest

Welcome to the **Go Bank Transfer Quest**! In this challenge, you will implement a `POST /v1/transfers` endpoint in a pre-scaffolded Go bank service. 

This quest focuses on idiomatic Go HTTP handler patterns, OpenTelemetry tracing, structured logging (`log/slog`), JWT authentication/authorization, and handler testing — without getting distracted by database or repository concerns. 

Everything below the API layer is pre-built. If you want to understand how the underlying layers work, check out the [Go Bank Architecture](../../bank/README.md).

The **account handler is your fully working reference** — read it, understand every pattern, and then replicate it for transfers.

## Your Quests

Work through the quests in order. Each step builds on the previous one.

### Quest 1: OpenAPI Spec Design

**File:** [docs/openapi/transfers.yaml](../../../docs/openapi/transfers.yaml)

**Context:**
Before writing code, we design the contract. Designing APIs contract-first ensures frontend and backend engineers agree on the API shape without waiting for the implementation.

**Task:**
Complete the partially filled `transfers.yaml` spec. 
- Define the request body schema. It needs `from_account_id` (string), `to_account_id` (string), and `amount` (integer).
- Define responses for success (200) returning `{ "status": "completed" }`.
- Define responses for various error scenarios (400, 403, 404, 422, 500), following the patterns established in `accounts.yaml`.

**Definition of Done:**
- Your `transfers.yaml` clearly maps out the endpoint, required properties, and all possible HTTP error codes.
- You can compare it side-by-side with `accounts.yaml` and see they share the same consistent structure.

### Quest 2: Wire the Routes

**File:** [internal/bank/api/server.go](../../bank/api/server.go)

**Context:**
The API server uses Gin as the HTTP router. It groups routes logically and applies middleware (like JWT authorization, logging, and tracing). The accounts group is already wired up as your reference.

**Task:**
- Open `internal/bank/api/server.go`.
- Uncomment the `transferHandler := transfer.New(svc)` initialization.
- Wire the `POST /v1/transfers` route using the exact same pattern as the accounts routes.
- Apply `middleware.JWTMiddleware` to the new group.
- Apply `middleware.RequireScope("transfers:write")` to the specific route.

**Definition of Done:**
- The Go code compiles successfully.
- You can run the following command with no errors:
  ```bash
  go build ./internal/bank/api/...
  ```

### Quest 3: Implement the Handler

**File:** [internal/bank/api/transfer/handler.go](../../bank/api/transfer/handler.go)

**Context:**
This is the core of the quest. You need to implement the `CreateTransfer` HTTP handler. You will extract the request body, start an OpenTelemetry trace, check business authorization rules (does the caller own the source account?), call the service layer, and accurately map domain errors to HTTP errors.

**Task:**
The handler skeleton contains 5 guided `TODO`s. Each `TODO` points directly to the exact line in the reference `api/account/handler.go` that demonstrates the pattern.
1. **Parse and Validate:** Use `c.ShouldBindJSON(&req)`. Return a 400 Bad Request on error using the `apierror` package.
2. **OpenTelemetry Trace:** Start a span using `otel.Tracer("bank").Start(ctx, "transfer.create")` and set the attributes for the request (`from_account_id`, `to_account_id`, `amount`). Ensure you `defer span.End()`.
3. **Verify Ownership:** Extract the caller's identity via `middleware.ClaimsFromCtx(ctx)`. Fetch the source account (`h.svc.GetAccount`). If the caller (`claims.Subject`) is not the account owner, return a 403 Forbidden.
4. **Call Service & Map Errors:** Call `h.svc.Transfer(...)`. Use the `errors.Is` switch pattern to map domain errors to HTTP errors:
   - `domain.ErrAccountNotFound` -> 404
   - `domain.ErrInsufficientFunds` -> 422 Unprocessable Entity
   - `domain.ErrAccountLocked` -> 422 Unprocessable Entity
   - `default` -> 500 Internal Server Error
5. **Log & Return:** Use `slog.InfoContext` to log the successful transfer (which automatically injects trace IDs). Return a 200 JSON response.

**Definition of Done:**
- Code compiles without syntax errors or unused variables.
- You can run the following command cleanly:
  ```bash
  go build ./internal/bank/api/transfer/...
  ```

### Quest 4: Write Handler Tests

**File:** [internal/bank/api/transfer/handler_test.go](../../bank/api/transfer/handler_test.go)

**Context:**
We use table-driven tests and `httptest` to unit test the handler isolated from the network. We mock the `Service` layer using `testify` (via `mockery`). 

**Task:**
- Open `internal/bank/api/transfer/handler_test.go`.
- Read how the happy path (200) and invalid body (400) cases are constructed.
- Implement the "wrong owner" (403) case: Mock `GetAccount` to return an account owned by "bob", but set the test token subject to "alice".
- Implement "insufficient funds" (422): Mock `GetAccount` successfully, but mock `Transfer` to return `domain.ErrInsufficientFunds`.
- Implement "source account not found" (404): Mock `GetAccount` to return `domain.ErrAccountNotFound`.

**Definition of Done:**
- Run the handler tests:
  ```bash
  go test ./internal/bank/api/transfer/... -v
  ```
- All tests pass successfully, confirming your handler perfectly maps edge cases.
- Finally, verify the entire service is still green:
  ```bash
  make test-bank
  ```

### Bonus Quest: Transfer Client & CLI

**Context:**
APIs are useless without clients. Building a strongly-typed Go client makes integration easy for CLI tools, web frontends, and other microservices.

**Task:**
1. **File:** [pkg/client/bank/client.go](../../../pkg/client/bank/client.go) — Implement the `Transfer` method. Look at `GetAccount` to see how we build the URL, set headers (especially `Authorization: Bearer`), serialize the JSON body, and execute `httppkg.DoRequest`.
2. **File:** [internal/bank/cli/transfer/transfer.go](../../bank/cli/transfer/transfer.go) — Wire up the CLI command. Parse the CLI arguments and invoke your newly written `bankClient.Transfer` method.

**Definition of Done:**
- Start the infrastructure:
  ```bash
  make db-up
  ```
- Build the binaries:
  ```bash
  make build
  ```
- Run the API server locally (in a separate terminal):
  ```bash
  make run-bank-api
  ```
- Issue a test JWT using `POST /v1/token`.
- Export the token so the CLI can use it:
  ```bash
  export BANK_TOKEN="<your-token-here>"
  ```
- Check out the CLI's built-in documentation:
  ```bash
  ./bin/bank-cli help
  ```
- Execute your CLI and see the balance successfully moved!
  ```bash
  ./bin/bank-cli transfer create ACC-1 ACC-2 5000
  ```

---
**Good luck! Remember to use the `account` handler as your ultimate reference guide.**