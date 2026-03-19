# Go Training Workshop

Welcome to the immersive Go Training Workshop! This program is designed to deliver a profound understanding of Go programming, transitioning from interpreted or higher-level languages to building production-grade services and command-line tools.

Throughout this hands-on workshop, delve into the nuances of building robust applications employing idiomatic design patterns, the standard library, and modern frameworks. Under expert guidance, unravel the design philosophies and engineering decisions that underpin effective Go development.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Getting Started](#getting-started)
3. [Module 1: Modern API Engineering Principles](#module-1-modern-api-engineering-principles)
4. [Module 2: Go Language Fundamentals](#module-2-go-language-fundamentals)
5. [Module 3: Building the Data & API Service](#module-3-building-the-data--api-service)
6. [Module 4: Temporal Orchestration](#module-4-temporal-orchestration)
7. [Challenges](#challenges)
8. [Further Learning & Resources](#further-learning--resources)

## Prerequisites

- [Go 1.26.1+](https://go.dev/dl/) installed.
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed (for database infrastructure).
- Basic experience with command-line tools is required.
- Familiarity with at least one other programming language (e.g., Python, Bash, or Java).

### Recommended Pre-Workshop Reading
To get the most out of this workshop, we recommend familiarising yourself with these foundational resources:
- [A Tour of Go (Basics 1)](https://go.dev/tour/basics/1) — Start here to get a feel for Go syntax and primitives.
- [Go Project Layout Standards](https://github.com/golang-standards/project-layout) — A popular convention for understanding workspace organisation.
- [Using Go Modules](https://go.dev/blog/using-go-modules) — Core dependency management in Go.
- [Go by Example](https://gobyexample.com/) — A hands-on introduction to Go using annotated example programs.
- [Effective Go](https://go.dev/doc/effective_go) — Tips for writing clear, idiomatic, and "Gopher-friendly" code.

## Getting Started

1. **Clone the Repository**:
    ```bash
    git clone https://github.com/romangurevitch/GoTraining.git
    cd GoTraining
    ```

2. **Makefile Help**:
    View the available make targets and their descriptions:
    ```bash
    make help
    ```

3. **Verify Installation**:
    Initialise dependencies and run the build/test suite:
    ```bash
    go mod tidy
    make build
    make test
    ```

4. **Start Infrastructure**:
    Launch the required database services:
    ```bash
    make db-up
    ```

5. **Open the Project in an IDE**:
    Two popular choices for Go development are:
    - [GoLand](https://www.jetbrains.com/go/): A powerful IDE by JetBrains dedicated to Go.
    - [Visual Studio Code (VSCode)](https://code.visualstudio.com/): Free editor with the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go).

## Module 1: Modern API Engineering Principles

Explore the foundational concepts for building production-ready APIs and platform tools:

- [API Fundamentals](docs/module1-fundamentals.md) — REST vs. RPC, Idempotency, and Contract-First design.
- [Fundamentals Overview](internal/fundamentals/README.md) — Security, Identity (AuthN/AuthZ), and Structured Logging (`slog`).

## Module 2: Go Language Fundamentals

Dive into the building blocks of Go by exploring the following topics:

- [The Mental Shift](internal/basics)
    - [Pointers](internal/basics/pointers/README.md)
    - [Type Assertions](internal/basics/casting/README.md)
    - [Parameters](internal/basics/parameters/README.md)
- [Structs & Layout](internal/basics)
    - [Entities](internal/basics/entity/README.md)
    - [Package Layout](internal/basics/layout/README.md)
    - [Embedding](internal/basics/embed/README.md)
- [Behaviours](internal/basics)
    - [Receivers](internal/basics/receivers/README.md)
    - [init()](internal/basics/init/README.md)
    - [Error Handling](internal/basics/err/README.md)
    - [Interfaces](internal/basics/interface/README.md)
- [Concurrency & Context](internal/basics)
    - [Concurrency](internal/basics/concurrency/README.md)
    - [Context](internal/basics/context/README.md)
- [Testing & Benchmarking](internal/basics)
    - [Testing](internal/basics/testing/README.md)
    - [Testify](internal/basics/testify/README.md)
    - [Benchmark](internal/basics/benchmark/README.md)
    - [HTTP Testing](internal/basics/httptest/README.md)
- [Advanced Features](internal/basics)
    - [Generics](internal/basics/generics/README.md)
    - [Mocking](internal/basics/mocking/README.md)
    - [Build Tags](internal/basics/buildtags/README.md)

Navigate to the respective [directories](internal/basics) to find code examples and documentation.

## Module 3: Building the Data & API Service

Build a persistent storage layer and HTTP service for the "Go Bank" domain. This module demonstrates a clean, layered architecture separating domain logic from transport and storage.

- Dive into the [Go Bank Architecture](internal/bank/README.md) to understand how the layers fit together.
- Ready to code? Jump straight into the [Go Bank Transfer Quest](internal/challenges/bank/README.md).

## Module 4: Temporal Orchestration

Discover reliable, durable execution patterns for long-running workflows:

- [Temporal Overview](internal/temporal/README.md) — Workflow vs. Activity and the Replay model.
- [Concepts Guide](docs/module4-temporal.md) — Why Temporal beats raw goroutines for distributed systems.
- [Worker Entrypoint](cmd/worker/main.go) — Temporal worker implementation stub.

## Challenges

Take on various exercises to test your understanding of Go:

- [Challenges Overview](internal/challenges/README.md)
- [Fix Me](internal/challenges/basics/fixme/README.md) — Diagnose and fix buggy code.
- [Implement Me](internal/challenges/basics/implme/README.md) — Complete the implementation to pass tests.
- [Go Bank Transfer Quest](internal/challenges/bank/README.md) — Build the `POST /v1/transfers` API endpoint.

Navigate to the respective [directories](internal/challenges) to start the exercises.

## Further Learning & Resources

Deepen your Go expertise with these advanced topics and ecosystem libraries:

### Testing & Quality
- [Testify](https://github.com/stretchr/testify) — A popular framework for writing expressive and robust tests.
- [Table-Driven Unit Tests in Go](https://dev.to/boncheff/table-driven-unit-tests-in-go-407b) — The idiomatic way to test multiple scenarios.
- [Roman's Presentation on Testing](https://youtu.be/ckzPKMhKNZk) — Expert insights into effective testing strategies.
- [Linting in Go (golangci-lint)](https://golangci-lint.run/) — Ensuring code consistency and catching common errors.

### API & Web Development
- [Structured Logging with slog](https://go.dev/blog/slog) — Native structured logging for production services.
- [Gin Web Framework](https://github.com/gin-gonic/gin) — High-performance HTTP web framework.
- [Introduction to gRPC](https://grpc.io/docs/what-is-grpc/) — High-performance, open source universal RPC framework.
- [Middlewares in Go](https://dev.to/karankumarshreds/middlewares-in-go-41j) — Patterns for handling cross-cutting concerns.

### Data & Configuration
- [go-jet](https://github.com/go-jet/jet) — Type-safe SQL builder with automatic scanning.
- [Cobra](https://github.com/spf13/cobra) — A library for creating powerful modern CLI applications.
- [Viper](https://github.com/spf13/viper) — Complete configuration solution including environment variables and config files.

### Concurrency & Performance
- [Go Concurrency Patterns](https://go.dev/blog/concurrency-timeouts) — Techniques for timeouts and robust concurrent code.
- [Benchmark Testing in Go](https://pkg.go.dev/testing#hdr-Benchmarks) — Measuring and optimising your application's performance.
- [Locks in Go](https://go.dev/tour/concurrency/9) — Understanding `sync.Mutex` and shared memory.
- [Temporal Documentation](https://docs.temporal.io/) — Durable execution for critical business workflows.
