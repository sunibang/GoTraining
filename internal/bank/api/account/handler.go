package account

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/service"
	"github.com/romangurevitch/go-training/pkg/api/apierror"
)

// Handler handles account-related HTTP requests.
type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// GetAccount handles GET /v1/accounts/:id
func (h *Handler) GetAccount(c *gin.Context) {
	ctx := c.Request.Context()

	ctx, span := otel.Tracer("bank").Start(ctx, "account.get")
	defer span.End()

	id := c.Param("id")
	slog.InfoContext(ctx, "get account", slog.String("account_id", id))

	result, err := h.svc.GetAccount(ctx, id)
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
	case err != nil:
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "could not get account", err))
	default:
		span.SetAttributes(attribute.String("account.owner", result.Owner))
		c.JSON(http.StatusOK, toAccountResponse(result))
	}
}

// CreateAccount handles POST /v1/accounts
func (h *Handler) CreateAccount(c *gin.Context) {
	ctx := c.Request.Context()

	ctx, span := otel.Tracer("bank").Start(ctx, "account.create")
	defer span.End()

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(apierror.NewAPIError(ctx, http.StatusBadRequest, "bad request", err))
		return
	}

	result, err := h.svc.CreateAccount(ctx, req.Owner)
	switch {
	case errors.Is(err, domain.ErrAccountAlreadyExists):
		c.JSON(apierror.NewAPIError(ctx, http.StatusConflict, "account already exists", err))
	case err != nil:
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "could not create account", err))
	default:
		span.SetAttributes(attribute.String("account.id", result.ID))
		slog.InfoContext(ctx, "account created", slog.String("account_id", result.ID))
		c.JSON(http.StatusCreated, toAccountResponse(result))
	}
}
