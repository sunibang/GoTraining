# Setup Guide

## Prerequisites

- [Go 1.26.1+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (for database infrastructure)

## Clone and Install

```bash
git clone https://github.com/romangurevitch/GoTraining.git
cd GoTraining
go mod tidy
```

## Verify Setup

```bash
make help    # view available commands
make build   # compiles all binaries
make test    # runs all tests
make lint    # lints the codebase
```

## Start Infrastructure

```bash
make db-up   # starts Postgres 15 via docker-compose
```

Postgres will be available at `localhost:5432` with:
- DB: `gobank`
- User: `gotrainer`
- Password: `verysecret`

## IDE Setup

**Visual Studio Code (VSCode):** Install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go).

**GoLand / IntelliJ:** No plugins required — Go is natively supported.
