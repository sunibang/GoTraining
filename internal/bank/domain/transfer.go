package domain

// Transfer represents a completed fund movement between two accounts.
// Used for logging and response mapping — not persisted as its own record.
// The underlying transactions are recorded as TypeDeposit and TypeWithdrawal entries.
type Transfer struct {
	FromAccountID string
	ToAccountID   string
	Amount        int64 // in cents
}
