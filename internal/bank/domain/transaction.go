package domain

import (
	"time"
)

// TransactionType defines whether money is entering or leaving an account.
type TransactionType string

const (
	TypeDeposit    TransactionType = "DEPOSIT"
	TypeWithdrawal TransactionType = "WITHDRAWAL"
)

// Transaction records a movement of funds.
type Transaction struct {
	ID        string          `json:"id"`
	AccountID string          `json:"account_id"`
	Amount    int64           `json:"amount"` // in cents
	Type      TransactionType `json:"type"`
	CreatedAt time.Time       `json:"created_at"`
}
