# API Layer

Welcome to the **API Layer**—the front door of Go Bank! This package handles all incoming HTTP requests using the [Gin](https://github.com/gin-gonic/gin) web framework. 

## What happens at the front door?

- **Routing:** Mapping URLs like `/v1/transfers` to the right Go functions.
- **Parsing & Validation:** Extracting JSON payloads and ensuring the data is correct before our service layer ever sees it.
- **Middleware Magic:** Orchestrating cross-cutting concerns like JWT authorisation, OpenTelemetry tracing, and structured logging (`slog`).
- **Response Formatting:** Returning clean JSON successes or translating domain errors into the right HTTP status codes.

## How to navigate

- **`account/`**: Your golden reference! This is a fully built, working implementation of GET/POST handlers, complete with tests. Read this first.
- **`transfer/`**: Your blank canvas. This is where you'll build the transfer endpoint.
- **`middleware/`**: Shared tools that wrap our handlers (auth, logging, tracing).
- **[server.go](server.go)**: The central hub that wires the router, middleware, and routes together.

## Ready to Code?

You've explored the engine room, from the inner domain to the outer API. This is where the magic happens. Your mission is to wire up the transfer routes, build the handler, and write the tests. 

Jump into the **[Go Bank Transfer Quest](../../challenges/bank/README.md)** to get started!