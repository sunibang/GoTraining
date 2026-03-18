package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	apierror "github.com/romangurevitch/go-training/pkg/api/error"
)

// claimsKey is an unexported named type used as a context key for JWT claims.
type claimsKey struct{}

// Claims extends jwt.RegisteredClaims with a Scope field.
// JWT payload example: { "sub": "alice", "scope": "accounts:read transfers:write", "exp": ... }
type Claims struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

// JWTMiddleware validates the Bearer token from Authorization header.
// On success: injects *Claims into the request context (c.Request.Context()) under claimsKey{};
// retrieve with ClaimsFromCtx(ctx).
// On failure: returns 401 and aborts — downstream handlers do not run.
func JWTMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(apierror.NewUnauthorizedError())
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(apierror.NewUnauthorizedError())
			c.Abort()
			return
		}

		// Store claims on request context so ClaimsFromCtx(ctx context.Context) works
		// in handlers and tests that receive ctx from c.Request.Context().
		ctx := context.WithValue(c.Request.Context(), claimsKey{}, claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RequireScope checks that the injected Claims contain the required scope.
// Returns 403 if scope is missing — used as per-route middleware after JWTMiddleware.
func RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFromCtx(c.Request.Context())
		if claims == nil || !hasScope(claims.Scope, scope) {
			c.JSON(http.StatusForbidden, &apierror.APIError{Message: "insufficient scope"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// hasScope checks whether the space-separated tokenScopes string contains the required scope token.
func hasScope(tokenScopes, required string) bool {
	for _, s := range strings.Fields(tokenScopes) {
		if s == required {
			return true
		}
	}
	return false
}

// ClaimsFromCtx extracts the JWT claims injected by JWTMiddleware from the request context.
// Returns nil if not present (should not happen in protected routes).
func ClaimsFromCtx(ctx context.Context) *Claims {
	if c, ok := ctx.Value(claimsKey{}).(*Claims); ok {
		return c
	}
	return nil
}
