package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/repository"
)

// Service is the business logic interface. Enables mock injection in handler tests.
type Service interface {
	CreateAccount(ctx context.Context, owner string) (*domain.Account, error)
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
	Deposit(ctx context.Context, accountID string, amount int64) error
	Withdraw(ctx context.Context, accountID string, amount int64) error
	Transfer(ctx context.Context, fromID, toID string, amount int64) error
}

// BankService implements Service backed by a Repository.
type BankService struct {
	repo repository.Repository
}

// Ensure BankService implements Service at compile time.
var _ Service = (*BankService)(nil)

func NewBankService(repo repository.Repository) *BankService {
	return &BankService{repo: repo}
}

func (s *BankService) CreateAccount(ctx context.Context, owner string) (*domain.Account, error) {
	acc := &domain.Account{
		ID:        fmt.Sprintf("ACC-%d", time.Now().UnixNano()),
		Owner:     owner,
		Balance:   0,
		Status:    domain.StatusOpen,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.SaveAccount(ctx, acc); err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}

	slog.InfoContext(ctx, "account created", slog.String("id", acc.ID), slog.String("owner", acc.Owner))
	return acc, nil
}

func (s *BankService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	return s.repo.GetAccount(ctx, id)
}

func (s *BankService) Deposit(ctx context.Context, accountID string, amount int64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	acc, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if err := acc.CanPerformTransaction(); err != nil {
		return err
	}

	acc.Balance += amount
	acc.UpdatedAt = time.Now()

	if err := s.repo.SaveAccount(ctx, acc); err != nil {
		return err
	}

	t := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d", time.Now().UnixNano()),
		AccountID: accountID,
		Amount:    amount,
		Type:      domain.TypeDeposit,
		CreatedAt: time.Now(),
	}

	return s.repo.SaveTransaction(ctx, t)
}

func (s *BankService) Withdraw(ctx context.Context, accountID string, amount int64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	acc, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if err := acc.CanPerformTransaction(); err != nil {
		return err
	}

	if acc.Balance < amount {
		return domain.ErrInsufficientFunds
	}

	acc.Balance -= amount
	acc.UpdatedAt = time.Now()

	if err := s.repo.SaveAccount(ctx, acc); err != nil {
		return err
	}

	t := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d", time.Now().UnixNano()),
		AccountID: accountID,
		Amount:    amount,
		Type:      domain.TypeWithdrawal,
		CreatedAt: time.Now(),
	}

	return s.repo.SaveTransaction(ctx, t)
}

// Transfer moves funds from one account to another.
// Pre-built for participants — they call this from the transfer handler.
func (s *BankService) Transfer(ctx context.Context, fromID, toID string, amount int64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	from, err := s.repo.GetAccount(ctx, fromID)
	if err != nil {
		return err
	}

	to, err := s.repo.GetAccount(ctx, toID)
	if err != nil {
		return err
	}

	if err := from.CanPerformTransaction(); err != nil {
		return err
	}

	if err := to.CanPerformTransaction(); err != nil {
		return err
	}

	if from.Balance < amount {
		return domain.ErrInsufficientFunds
	}

	from.Balance -= amount
	from.UpdatedAt = time.Now()

	to.Balance += amount
	to.UpdatedAt = time.Now()

	if err := s.repo.SaveAccount(ctx, from); err != nil {
		return fmt.Errorf("failed to debit source account: %w", err)
	}

	if err := s.repo.SaveAccount(ctx, to); err != nil {
		return fmt.Errorf("failed to credit destination account: %w", err)
	}

	debit := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d-D", time.Now().UnixNano()),
		AccountID: fromID,
		Amount:    amount,
		Type:      domain.TypeWithdrawal,
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveTransaction(ctx, debit); err != nil {
		return err
	}

	credit := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d-C", time.Now().UnixNano()),
		AccountID: toID,
		Amount:    amount,
		Type:      domain.TypeDeposit,
		CreatedAt: time.Now(),
	}

	slog.InfoContext(ctx, "transfer completed",
		slog.String("from_account_id", fromID),
		slog.String("to_account_id", toID),
		slog.Int64("amount", amount),
	)

	return s.repo.SaveTransaction(ctx, credit)
}
