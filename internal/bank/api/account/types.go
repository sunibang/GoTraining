package account

import (
	"time"

	"github.com/romangurevitch/go-training/internal/bank/domain"
)

// CreateAccountRequest is the JSON body for POST /v1/accounts.
type CreateAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
}

// AccountResponse is the JSON body returned by GET /v1/accounts/:id and POST /v1/accounts.
type AccountResponse struct {
	ID        string    `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toAccountResponse(a *domain.Account) AccountResponse {
	return AccountResponse{
		ID:        a.ID,
		Owner:     a.Owner,
		Balance:   a.Balance,
		Status:    string(a.Status),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
