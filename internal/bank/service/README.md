# Service Layer

Welcome to the **Service Layer**—the brain of the Go Bank! This package sits right between our HTTP API and the database repository, orchestrating the actual banking operations.

## The Core Rules

- **Enforce Business Logic:** We make sure you can't withdraw more money than you have, or use a locked account.
- **Orchestrate Multi-Step Moves:** A money transfer isn't just one database query! The service coordinates deducting from one account and crediting another safely.
- **Speak in Domain Terms:** The service returns our specific `domain` errors (like `ErrInsufficientFunds`), letting the API layer figure out the right HTTP status code.

## Built for Testing

We interact with this layer through the [Service](service.go) interface. Why? So we can easily mock it out when testing our API handlers! (You can find the generated mock in [mocks/mock_service.go](mocks/mock_service.go)).

## Your Next Step

This layer is completely ready to go. 

Read through the `Transfer` method in [service.go](service.go) to see how it enforces rules and uses the repository. The final step before building is to see how the HTTP requests come in. Head over to the **[API Layer](../api/README.md)**!