package api

import "time"

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

// CreateTransferRequest is the JSON body for POST /v1/transfers.
type CreateTransferRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required"`
	ToAccountID   string `json:"to_account_id"   binding:"required"`
	Amount        int64  `json:"amount"           binding:"required,gte=1"`
}

// TransferResponse is the JSON body returned on successful transfer.
type TransferResponse struct {
	Status string `json:"status"`
}
