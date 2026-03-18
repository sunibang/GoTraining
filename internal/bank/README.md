# Go Bank Architecture

Welcome to the engine room! This is the core service implementation for the "Go Bank" domain. We've structured this service using a classic layered architecture to keep our code clean, testable, and easy to maintain.

## Explore the Layers

Each layer has a specific job. Click through to explore how they work:

- **[Domain](domain/README.md)**: The heart of our business. Contains core entities like [Account](domain/account.go) and [Transaction](domain/transaction.go). It depends on nothing else.
- **[Repository](repository/README.md)**: The data vault. Handles all Postgres database interactions using type-safe SQL queries via `go-jet`.
- **[Service](service/README.md)**: The brain. Orchestrates operations and enforces business rules (like making sure you have enough funds before a transfer).
- **[API](api/README.md)**: The front door. Exposes REST endpoints using Gin, complete with structured logging, OpenTelemetry tracing, and JWT auth.

*Curious how it all boots up? The `app` and `config` packages wire these layers together in [cmd/bank-api/main.go](../../cmd/bank-api/main.go).*

## Your Exploration Journey

Before diving into the code, let's explore the architecture from the inside out. 

Start your journey at the absolute center: the **[Domain Layer](domain/README.md)**.