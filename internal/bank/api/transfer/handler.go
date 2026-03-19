package transfer

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/service"
	"github.com/romangurevitch/go-training/pkg/api/apierror"
)

// Handler handles transfer-related HTTP requests.
type Handler struct {
	svc service.Service
}

// New creates a new transfer handler.
func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// CreateTransfer handles POST /v1/transfers.
func (h *Handler) CreateTransfer(c *gin.Context) {
	ctx := c.Request.Context()

	// TODO 1: Parse and validate request body.

	// TODO 2: Start an OTel span and set attributes.

	// TODO 3: Verify ownership.

	// TODO 4: Call service and map errors.

	// TODO 5: Log success and return 200.

	// REMOVE BELOW LINES when you implement TODO 1:
	_, _ = errors.New(""), middleware.ClaimsFromCtx(ctx) // silence unused imports
	_ = domain.ErrAccountNotFound                        // silence unused import
	_ = apierror.ErrInternalServerError                  // silence unused import
}
