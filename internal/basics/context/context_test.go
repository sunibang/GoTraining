package context

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 1. Basic Cancellation
// Demonstrate how to manually trigger cancellation.
func TestContextWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel() // Manually trigger cancellation
	}()

	select {
	case <-ctx.Done():
		assert.Equal(t, context.Canceled, ctx.Err())
	case <-time.After(1 * time.Second):
		t.Fatal("Context was never cancelled")
	}
}

// 2. Timeout (The most common usage)
// Automatically cancels after a duration.
func TestContextWithTimeout(t *testing.T) {
	// A context that expires in 50ms
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel() // Always call cancel, even on timeout

	select {
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for context timeout!")
	case <-ctx.Done():
		// Context should be cancelled by the timer
		assert.Equal(t, context.DeadlineExceeded, ctx.Err())
	}
}

// 3. New in Go 1.24: t.Context()
// This context is automatically cancelled when the test finishes.
func TestTestingContext(t *testing.T) {
	ctx := t.Context()

	// Use the context in a long-running operation
	res, err := doWork(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 42, res)

	// No need to call cancel() - Go handles it!
}

// 4. Use Case: Authentication
type User struct {
	ID   string
	Role string
}

type authKey int
const userKey authKey = 0

func TestContextForAuth(t *testing.T) {
	user := User{ID: "user-123", Role: "admin"}

	// Set user in context (e.g., in a middleware)
	ctx := context.WithValue(context.Background(), userKey, user)

	// Retrieve user from context (e.g., in a handler/service)
	val := ctx.Value(userKey)
	assert.NotNil(t, val)

	retrievedUser, ok := val.(User)
	assert.True(t, ok)
	assert.Equal(t, "admin", retrievedUser.Role)
}

// 5. Use Case: Tracing
type traceKey string
const traceIDKey traceKey = "trace-id"

func TestContextForTracing(t *testing.T) {
	ctx := context.WithValue(context.Background(), traceIDKey, "tx-999-abc")

	// In a deeply nested function:
	traceID := ctx.Value(traceIDKey).(string)
	assert.Equal(t, "tx-999-abc", traceID)
}

// 6. Practical Pattern: Work Function
// A common way to structure functions that respect context.
func doWork(ctx context.Context) (int, error) {
	// Simulate some async work
	resCh := make(chan int)
	go func() {
		time.Sleep(50 * time.Millisecond)
		resCh <- 42
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err() // Return the reason for cancellation
	case res := <-resCh:
		return res, nil
	}
}

func TestDoWorkRespectsContext(t *testing.T) {
	t.Run("SuccessCase", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		val, err := doWork(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	t.Run("TimeoutCase", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		val, err := doWork(ctx)
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Equal(t, 0, val)
	})
}
