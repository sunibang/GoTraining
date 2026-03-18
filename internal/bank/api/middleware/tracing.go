package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// TracingMiddleware starts an OTel span for every incoming HTTP request.
// Span name: "{method} {path}" — e.g. "POST /v1/transfers"
// Uses otelgin which handles W3C trace context propagation automatically.
func TracingMiddleware(serviceName string) gin.HandlerFunc {
	// Fallback keeps spans identifiable in Jaeger if the caller omits a service name.
	if serviceName == "" {
		serviceName = "bank-api"
	}
	return otelgin.Middleware(serviceName)
}
