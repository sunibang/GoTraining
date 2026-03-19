# Challenges

This directory contains all student exercises for the Go Training workshop.

## Structure

```
challenges/
├── basics/
│   ├── fixme/    # Find and fix the bugs
│   └── implme/   # Implement the function
└── bank/         # Go Bank service quests
```

## basics/fixme

Short, focused exercises where buggy code is provided. Your task: identify the problem and fix it.
Inspired by the [ConcurrencyWorkshop](https://github.com/) fixme pattern.

## basics/implme

Exercises with `panic("implement me!")` stubs. Your task: implement the function to make the tests pass.

## bank/

The **[Go Bank Transfer Quest](bank/README.md)** is our main challenge! 

You'll implement the `POST /v1/transfers` API endpoint in a pre-scaffolded service, focusing on:
- Idiomatic HTTP handler patterns using Gin.
- OpenTelemetry tracing and structured logging with `slog`.
- JWT authentication and scope-based authorisation.
- Table-driven unit testing for handlers.

Everything below the API layer is pre-built so you can focus on building production-grade APIs.

Run tests with: `make test-bank`

