# Challenge 05: Build a Safe HTTP Client

**Difficulty:** Basic
**Covers:** `http`, `httptest`, `pointers`

---

## Goal

Build a reusable HTTP client package that handles the most common pitfalls in production Go code.

---

## Requirements

1. **Configurable client** — create a `Client` struct with a configurable timeout. Use the functional options pattern so callers can opt into non-default settings without breaking the API.

2. **Safe `Get` helper** — implement `Get(url string) ([]byte, error)` that:
   - Always closes the response body (even on error).
   - Returns an error for non-2xx status codes (Go's `http.Client` does **not** do this automatically).
   - Returns the response body bytes on success.

3. **Tests** — write `httptest`-based tests covering:
   - A successful 200 response returning the expected body.
   - A 404 response returning an error from `Get`.
   - A 500 response returning an error from `Get`.

---

## Starter Code

Complete the TODOs in `client.go` and `client_test.go`.

```go
// client.go
package safeclient

import (
    "fmt"
    "io"
    "net/http"
    "time"
)

type Client struct {
    http    *http.Client
    timeout time.Duration
}

type Option func(*Client)

func WithTimeout(d time.Duration) Option {
    // TODO: return an Option that sets c.timeout and updates c.http.Timeout
    panic("not implemented")
}

func New(opts ...Option) *Client {
    // TODO: create a Client with a 10s default timeout, then apply opts
    panic("not implemented")
}

// Get fetches url and returns the body bytes.
// Returns an error if the request fails OR if the status code is not 2xx.
func (c *Client) Get(url string) ([]byte, error) {
    // TODO:
    // 1. Use c.http.Get(url)
    // 2. Always defer resp.Body.Close()
    // 3. Check resp.StatusCode — return a descriptive error for non-2xx
    // 4. Read and return the body with io.ReadAll
    panic("not implemented")
}
```

---

## Skills Practiced

- HTTP client configuration
- Response body lifecycle (`defer Body.Close()`)
- Functional options pattern
- Testing with `httptest.NewServer`

---

## Hints

- `http.StatusOK` is 200. Any status `< 200` or `>= 300` is non-success.
- `fmt.Errorf("unexpected status %d", resp.StatusCode)` is a good error format.
- `io.ReadAll(resp.Body)` reads the full body into a `[]byte`.
- Always close the body *before* returning an error, otherwise you leak the connection.
