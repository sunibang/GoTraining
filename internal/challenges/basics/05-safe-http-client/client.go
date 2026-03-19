// Package safeclient is a challenge skeleton. Complete the TODOs.
package safeclient

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client wraps net/http.Client with safe defaults.
type Client struct {
	http    *http.Client
	timeout time.Duration
}

// Option is a functional option for Client.
type Option func(*Client)

// WithTimeout sets a custom timeout on the client.
func WithTimeout(d time.Duration) Option {
	// TODO: return an Option that sets c.timeout and updates c.http.Timeout
	panic("not implemented")
}

// New returns a Client with a 10-second default timeout.
// Additional options are applied after the defaults.
func New(opts ...Option) *Client {
	// TODO: create a Client with a 10s default timeout, then apply opts
	panic("not implemented")
}

// Get fetches url and returns the response body bytes.
// Returns an error if the request fails OR if the status code is not 2xx.
// The response body is always closed before returning.
func (c *Client) Get(url string) ([]byte, error) {
	// TODO:
	// 1. Use c.http.Get(url) to make the request
	// 2. Always defer resp.Body.Close()
	// 3. Check resp.StatusCode — return an error for non-2xx
	// 4. Read and return the body with io.ReadAll
	panic("not implemented")
}

// nonSuccessError returns a descriptive error for non-2xx responses.
//
//nolint:unused // challenge skeleton: students call this from Get once implemented
func nonSuccessError(statusCode int) error {
	return fmt.Errorf("unexpected status %d: %s", statusCode, http.StatusText(statusCode))
}

// ensure io is used (remove once you use io.ReadAll in Get)
var _ = io.ReadAll
