# Challenge 5: Embeddable Middleware Chain

**Difficulty:** Advanced
**Covers:** `interface`, `http`, `context`, `receivers`, `mocking`

---

## Goal

Build a composable HTTP middleware system from scratch. Each middleware wraps the next handler, forming a chain. This is the same pattern used by popular Go frameworks like Chi, Gin, and Echo.

---

## Requirements

### 1. Core types

```go
// Middleware wraps an http.Handler and returns a new http.Handler.
// Define this on the consumer side (in your package, not in net/http).
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares in order: Chain(h, A, B, C) → A(B(C(h)))
func Chain(h http.Handler, middlewares ...Middleware) http.Handler
```

### 2. Implement three middlewares

| Middleware | Behaviour |
|------------|-----------|
| `LoggingMiddleware` | Logs `method path` for every request |
| `AuthMiddleware` | Reads `X-Auth-Token` header; returns 401 if missing or not `"valid-token"` |
| `RateLimitMiddleware` | Allows at most N requests per second (use a simple counter + time window) |

### 3. Integration tests (httptest)

Using `httptest.NewServer`, write tests for:
- The full chain passes a valid request to the handler
- `AuthMiddleware` rejects requests with missing/wrong token (401)
- Removing `LoggingMiddleware` does not affect response correctness
- `RateLimitMiddleware` returns 429 after N+1 requests in the window

### 4. Unit tests with mocked handlers

Test each middleware in isolation using a mock `http.Handler`:
- `LoggingMiddleware`: verify the inner handler is called exactly once
- `AuthMiddleware`: verify inner handler is **not** called on 401
- Use `httptest.NewRecorder` as the response writer

### 5. (Optional) Benchmark

```go
func BenchmarkChain_1(b *testing.B)  { benchChain(b, 1) }
func BenchmarkChain_5(b *testing.B)  { benchChain(b, 5) }
func BenchmarkChain_10(b *testing.B) { benchChain(b, 10) }
```

---

## Skills Practiced

- Interface composition via function types
- Context values for request-scoped data
- HTTP handler chaining
- Testing with mocked `http.Handler`
- Benchmark methodology

---

## Hints

- `http.HandlerFunc` adapts a function to `http.Handler` — your Chain result can use it
- For auth: `r.Header.Get("X-Auth-Token")`
- For rate limiting, a simple approach: `sync/atomic` counter reset by a ticker
- To mock an `http.Handler`, implement the single `ServeHTTP(w, r)` method or use `http.HandlerFunc(func(w,r){...})`
