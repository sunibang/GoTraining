# Module 4: Temporal Orchestration

## Why Temporal?

Raw goroutines and channels are great for short-lived concurrent work, but break down for long-running operations that must survive failures, retries, and restarts. Temporal provides **durable execution**: workflows that automatically retry, resume, and track state across failures.

## Key Concepts

| Concept | Description |
|---|---|
| **Workflow** | Deterministic, replayable business logic. No I/O directly — delegates to Activities. |
| **Activity** | A single, retriable step (e.g., debit an account, send an email). Can have timeouts and retries. |
| **Worker** | A process that polls Temporal for tasks and executes workflows/activities. |
| **Task Queue** | The named queue a worker listens on. Workflows are dispatched to workers via task queues. |

## Order Processing Demo

The demo features two order processing workflows demonstrating durable execution, signal handling, and child workflows.

### 1. Start Services

Start the Temporal server and WireMock (Inventory API):

```bash
make temporal-up
```

- **Temporal Web UI**: http://localhost:8233
- **WireMock (Inventory API)**: http://localhost:8080

### 2. Start the Worker

The worker listens for tasks on the `order-processing` queue:

```bash
make worker-start
```

### 3. Execute a Workflow

Use the client to start either of the available workflows.

#### AutoProcessOrder (No signals required)

Automatically drives the order through every stage: `PLACED` → `PICKED` → `SHIPPED` → `COMPLETED`.

```bash
go run cmd/client/main.go \
  -workflow=AutoProcessOrder \
  -order='{
    "id": "00000000-0000-0000-0000-000000000001",
    "line_items": [
      {
        "product_id": "00000000-0000-0000-0000-000000000001",
        "quantity": 10,
        "price_per_item": "29.99"
      }
    ]
  }'
```

#### ProcessOrder (Signal-driven)

Pauses at each stage and waits for an external signal before continuing:

```bash
go run cmd/client/main.go \
  -workflow=ProcessOrder \
  -order='{
    "id": "00000000-0000-0000-0000-000000000001",
    "line_items": [
      {
        "product_id": "00000000-0000-0000-0000-000000000001",
        "quantity": 10,
        "price_per_item": "29.99"
      }
    ]
  }'
```

The workflow will:
1. Validate the order and check inventory.
2. Wait for a `pickOrder` signal (or `cancelOrder`).
3. Process payment via a child workflow.
4. Wait for `shipOrder` signal.
5. Wait for `markOrderAsDelivered` signal.

### Interacting with Signal-Driven Workflows

Send signals using the Temporal CLI or Web UI:

```bash
# Pick the order (moves from PLACED to PICKED)
temporal workflow signal --workflow-id order-<uuid> --name pickOrder

# Ship the order (moves to SHIPPED)
temporal workflow signal --workflow-id order-<uuid> --name shipOrder

# Mark as delivered (moves to COMPLETED)
temporal workflow signal --workflow-id order-<uuid> --name markOrderAsDelivered

# Or cancel the order (before picking)
temporal workflow signal --workflow-id order-<uuid> --name cancelOrder
```

Query the current status:

```bash
temporal workflow query --workflow-id order-<uuid> --name GetOrderStatus
```

## Self-Paced Resources

- [Temporal Go SDK documentation](https://docs.temporal.io/develop/go)
- [Temporal tutorials](https://learn.temporal.io/)
- [Saga pattern explained](https://microservices.io/patterns/data/saga.html)
