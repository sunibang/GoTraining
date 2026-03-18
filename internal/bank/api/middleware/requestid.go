package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

const headerRequestID = "X-Request-Id"

// RequestIDMiddleware accepts the X-Request-Id header if it contains a valid UUID.
// If the header is absent or contains an invalid UUID, a fresh UUID is generated.
// The request ID is injected into both the Gin context and c.Request.Context()
// so it is accessible via RequestIDFromCtx(ctx) in downstream handlers.
// The value is also echoed in the response header and read by slog-gin when WithRequestID: true.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(headerRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		} else if _, err := uuid.Parse(requestID); err != nil {
			requestID = uuid.New().String()
		}
		c.Writer.Header().Set(headerRequestID, requestID)
		c.Set(string(requestIDKey), requestID)
		ctx := context.WithValue(c.Request.Context(), requestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RequestIDFromCtx extracts the request ID from the Gin context.
func RequestIDFromCtx(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
