# Mocking in Go

## Description

Interfaces in Go are used to decouple packages. When unit testing code that depends on interfaces, mocking allows us to isolate the unit under test by providing controlled implementations of those dependencies.

In this project, we use two of the most popular mocking tools in the Go ecosystem.

## Mockery

[Mockery](https://github.com/vektra/mockery) is a tool that generates mocks for Golang interfaces using the [testify](https://github.com/stretchr/testify) framework.

### Pros

- Integrates seamlessly with `testify/mock`, the most popular assertion library in Go.
- Provides a powerful CLI for bulk mock generation.
- **Type-Safe Expectations:** Modern Mockery (v2+) supports a type-safe `EXPECT()` API similar to `gomock`, catching signature changes at compile time.
- Flexible argument matching (e.g., `mock.Anything`, `mock.MatchedBy`).

### Cons

- Requires `testify` as a dependency.

## Go Mock (uber/mock)

[GoMock](https://github.com/uber-go/mock) is the official fork of the now-archived `github.com/golang/mock`. It is a reflection-based mocking framework.

### Pros

- **Strict Type-Safety:** Mocks are generated as Go code that strictly follows the interface, ensuring compile-time correctness.
- **Call Ordering:** Excellent support for verifying the exact order of method calls.
- Part of the `uber-go` suite, known for high-quality engineering standards.

### Cons

- Slightly more verbose setup compared to Mockery/Testify.
- Less expressive assertions compared to Testify's `assert` and `require`.

## Which one should I use?

- Use **Mockery** if you are already using `testify` for assertions and want a consistent, expressive mocking experience.
- Use **GoMock** if you require strict call ordering or prefer the Uber-style engineering patterns.

In this workshop, we primarily use **Mockery** for the Go Bank service to keep our tests concise and readable.
