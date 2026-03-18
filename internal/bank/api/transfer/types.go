package transfer

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
