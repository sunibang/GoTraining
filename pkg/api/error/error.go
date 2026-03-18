package error

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// ErrInternalServerError is the sentinel for unexpected server errors.
var ErrInternalServerError = errors.New("internal server error")

// NewAPIError logs the error with context (trace_id/span_id injected automatically by OtelHandler)
// and returns a tuple suitable for c.JSON(apierror.NewAPIError(...)).
func NewAPIError(ctx context.Context, status int, msg string, err error) (int, *APIError) {
	slog.ErrorContext(ctx, msg,
		slog.Int("status", status),
		slog.Any("error", err),
	)
	return status, &APIError{Message: msg}
}

// NewUnauthorizedError returns a 401 without logging (avoids noise from probing).
func NewUnauthorizedError() (int, *APIError) {
	return http.StatusUnauthorized, &APIError{Message: "unauthorized"}
}

// APIError is the JSON error body returned by all API endpoints.
type APIError struct {
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %v", e.Message)
}
