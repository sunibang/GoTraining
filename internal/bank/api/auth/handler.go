package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
)

// Handler issues JWT tokens for the /v1/token endpoint.
// Pre-built — participants use curl to get tokens for manual testing.
type Handler struct {
	secret string
}

func New(secret string) *Handler {
	return &Handler{secret: secret}
}

// IssueTokenRequest is the request body for POST /v1/token.
type IssueTokenRequest struct {
	Sub   string `json:"sub"   binding:"required"` // account owner — becomes JWT Subject
	Scope string `json:"scope" binding:"required"` // space-separated scopes
}

// IssueTokenResponse wraps the signed token.
type IssueTokenResponse struct {
	Token string `json:"token"`
}

// IssueToken issues a signed JWT for the given sub and scope.
// No authentication required — this is a training tool, not production auth.
//
// Example request:
//
//	curl -X POST localhost:8080/v1/token \
//	  -H 'Content-Type: application/json' \
//	  -d '{"sub":"alice","scope":"accounts:read transfers:write"}'
func (h *Handler) IssueToken(c *gin.Context) {
	var req IssueTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request: " + err.Error()})
		return
	}

	claims := middleware.Claims{
		Scope: req.Scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   req.Sub,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, IssueTokenResponse{Token: signed})
}
