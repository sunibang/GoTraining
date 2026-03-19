# 🌐 Context in Go (`context`)

The `context` package is one of Go's most important tools for managing **cancellation**, **deadlines**, and **request-scoped values** across API boundaries and between goroutines.

---

## 1. What is Context?

A `Context` is a thread-safe object that carries information through a call graph. Its primary purpose is to signal when a process should stop working and return.

### 🖼️ The Context Tree
Contexts are hierarchical. When a parent context is cancelled, all contexts derived from it are also cancelled.

```text
       context.Background() (Root)
               |
        +------+------+
        |             |
   WithCancel    WithTimeout (5s)
        |             |
   WithValue     WithCancel
```

---

## 2. Core Use Cases

| Category | Typical Data/Action |
|----------|---------------------|
| **Auth** | User ID, Roles, Permissions (via `WithValue`) |
| **Tracing** | Trace ID, Span ID for observability (via `WithValue`) |
| **Concurrency** | Stopping background workers, limiting wait times |
| **I/O Management** | Aborting database queries or HTTP requests when a user disconnects |

---

## 3. Creating Contexts

| Function | Purpose |
|----------|---------|
| `context.Background()` | The root context; typically used in `main` or top-level requests. |
| `context.TODO()` | Use when you're unsure which context to use or it's not yet available. |
| `context.WithCancel(parent)` | Returns a child and a `cancel` function to stop it manually. |
| `context.WithTimeout(parent, duration)` | Automatically cancels after a specific duration. |
| `context.WithDeadline(parent, time)` | Automatically cancels at a specific clock time. |
| `context.WithValue(parent, key, val)` | Carries request-scoped data (use sparingly!). |

---

## 4. New in Go 1.24+: `t.Context()`

The `testing` package now provides a built-in context that is automatically cancelled when a test (and all its subtests) finishes. This is the **preferred** way to handle context in modern Go tests.

```go
func TestMyTask(t *testing.T) {
    ctx := t.Context() // Automatically cancelled when test finishes
    err := DoSomething(ctx)
    assert.NoError(t, err)
}
```

---

## 5. Best Practices (The "Golden Rules")

1.  **Pass as First Argument**: Context should always be the first parameter of a function: `func DoWork(ctx context.Context, ...)`.
2.  **Don't Store in Structs**: Never store a Context inside a struct; pass it explicitly to methods instead.
3.  **Always Call Cancel**: When using `WithCancel`, `WithTimeout`, or `WithDeadline`, always `defer cancel()`. This releases resources even if the work finishes early (not needed for `t.Context()`).
4.  **Context is Immutable**: You never "change" a context; you derive a new one from a parent.
5.  **Values for Metadata Only**: Use `WithValue` only for request-scoped data (e.g., trace IDs, auth tokens), not for passing optional parameters to functions.

---

## ⚠️ Critical Pitfall: Goroutine Leaks

If you create a `WithTimeout` context and don't call `cancel()`, the timer will keep running in the background until it expires, even if your function has already returned. This is a "goroutine leak."

```go
// GOOD
ctx, cancel := context.WithTimeout(parent, time.Hour)
defer cancel() // Timer is stopped when function returns
```

---

## 🧪 Running the Examples

Explore `context_test.go` for practical examples of cancellation, tracing, and `t.Context()`.

```bash
go test -v ./internal/basics/context/...
```
