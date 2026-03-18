package domain

import (
	"errors"
	"time"
)

// Business errors
var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrAccountLocked        = errors.New("account is locked")
	ErrInvalidAmount        = errors.New("invalid amount")
	ErrAccountAlreadyExists = errors.New("account already exists")
)

// AccountStatus represents the lifecycle state of an account.
type AccountStatus string

const (
	StatusOpen   AccountStatus = "OPEN"
	StatusLocked AccountStatus = "LOCKED"
	StatusClosed AccountStatus = "CLOSED"
)

// Account represents a bank account in the system.
type Account struct {
	ID        string        `json:"id"`
	Owner     string        `json:"owner"`
	Balance   int64         `json:"balance"` // in cents
	Status    AccountStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// CanPerformTransaction checks if the account is in a state that allows transactions.
func (a *Account) CanPerformTransaction() error {
	if a.Status == StatusLocked {
		return ErrAccountLocked
	}
	if a.Status == StatusClosed {
		return errors.New("account is closed")
	}
	return nil
}
