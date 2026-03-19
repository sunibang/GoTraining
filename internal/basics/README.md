# Module 2: Go Language Fundamentals

This module covers the core building blocks of the Go programming language. Each directory contains code examples and a focused `README.md` to guide you through the concepts.

## Topics

### The Mental Shift
- **[Pointers](pointers/README.md)** — Memory addresses and indirect values
- **[Type Assertions & Casting](casting/README.md)** — Working with dynamic types and interfaces
- **[Parameters](parameters/README.md)** — Passing values vs. pointers

### Structs & Layout
- **[Entities](entity/README.md)** — Defining data structures
- **[Package Layout](layout/README.md)** — Organizing your code idiomaticlly
- **[Embedding](embed/README.md)** — Composition over inheritance

### Behaviours
- **[Receivers](receivers/README.md)** — Adding methods to types
- **[init()](init/README.md)** — Package initialization
- **[Error Handling](err/README.md)** — Sentinel errors, wrapping, and the `error` interface
- **[Interfaces](interface/README.md)** — Implicit implementation and decoupling

### Concurrency & Context
- **[Concurrency](concurrency/README.md)** — Goroutines and Channels
- **[Context](context/README.md)** — Cancellation, deadlines, and request-scoped values

### Testing & Benchmarking
- **[Testing](testing/README.md)** — Unit testing and table-driven tests
- **[Testify](testify/README.md)** — Fluent assertions and requirements
- **[Benchmark](benchmark/README.md)** — Performance measurement
- **[HTTP Testing](httptest/README.md)** — Testing handlers without a network

### Advanced Features
- **[Generics](generics/README.md)** — Writing type-agnostic code
- **[Mocking](mocking/README.md)** — Using Mockery for dependency isolation
- **[Build Tags](buildtags/README.md)** — Conditional compilation
