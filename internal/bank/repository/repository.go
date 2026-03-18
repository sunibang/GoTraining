package repository

import (
	"context"

	"github.com/romangurevitch/go-training/internal/bank/domain"
)

// Repository defines the data access contract for the bank.
type Repository interface {
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
	SaveAccount(ctx context.Context, account *domain.Account) error
	ListTransactions(ctx context.Context, accountID string) ([]domain.Transaction, error)
	SaveTransaction(ctx context.Context, transaction *domain.Transaction) error
}
