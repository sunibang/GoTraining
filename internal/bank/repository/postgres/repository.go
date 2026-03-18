package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/repository"
	genmodel "github.com/romangurevitch/go-training/internal/bank/repository/postgres/gen/gobank/public/model"
	gentable "github.com/romangurevitch/go-training/internal/bank/repository/postgres/gen/gobank/public/table"
)

// Ensure PostgresRepository implements Repository at compile time.
var _ repository.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	stmt := gentable.Accounts.SELECT(gentable.Accounts.AllColumns).
		WHERE(gentable.Accounts.ID.EQ(postgres.String(id)))

	var dest genmodel.Accounts
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}
	return toDomainAccount(&dest), nil
}

func (r *PostgresRepository) SaveAccount(ctx context.Context, account *domain.Account) error {
	m := fromDomainAccount(account)
	stmt := gentable.Accounts.INSERT(
		gentable.Accounts.ID,
		gentable.Accounts.Owner,
		gentable.Accounts.Balance,
		gentable.Accounts.Status,
		gentable.Accounts.CreatedAt,
		gentable.Accounts.UpdatedAt,
	).VALUES(
		m.ID, m.Owner, m.Balance, m.Status, m.CreatedAt, m.UpdatedAt,
	).ON_CONFLICT(gentable.Accounts.ID).DO_UPDATE(
		postgres.SET(
			gentable.Accounts.Balance.SET(postgres.Int(m.Balance)),
			gentable.Accounts.Status.SET(postgres.String(m.Status)),
			gentable.Accounts.UpdatedAt.SET(postgres.TimestampzT(m.UpdatedAt)),
		),
	)
	_, err := stmt.ExecContext(ctx, r.db)
	return err
}

func (r *PostgresRepository) ListTransactions(ctx context.Context, accountID string) ([]domain.Transaction, error) {
	stmt := gentable.Transactions.SELECT(gentable.Transactions.AllColumns).
		WHERE(gentable.Transactions.AccountID.EQ(postgres.String(accountID))).
		ORDER_BY(gentable.Transactions.CreatedAt.DESC())

	var dest []genmodel.Transactions
	if err := stmt.QueryContext(ctx, r.db, &dest); err != nil {
		return nil, err
	}

	txs := make([]domain.Transaction, len(dest))
	for i, t := range dest {
		txs[i] = toDomainTransaction(&t)
	}
	return txs, nil
}

func (r *PostgresRepository) SaveTransaction(ctx context.Context, t *domain.Transaction) error {
	m := fromDomainTransaction(t)
	stmt := gentable.Transactions.INSERT(
		gentable.Transactions.ID,
		gentable.Transactions.AccountID,
		gentable.Transactions.Amount,
		gentable.Transactions.Type,
		gentable.Transactions.CreatedAt,
	).VALUES(m.ID, m.AccountID, m.Amount, m.Type, m.CreatedAt)
	_, err := stmt.ExecContext(ctx, r.db)
	return err
}

func toDomainAccount(m *genmodel.Accounts) *domain.Account {
	return &domain.Account{
		ID:        m.ID,
		Owner:     m.Owner,
		Balance:   m.Balance,
		Status:    domain.AccountStatus(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func fromDomainAccount(a *domain.Account) *genmodel.Accounts {
	return &genmodel.Accounts{
		ID:        a.ID,
		Owner:     a.Owner,
		Balance:   a.Balance,
		Status:    string(a.Status),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func toDomainTransaction(m *genmodel.Transactions) domain.Transaction {
	return domain.Transaction{
		ID:        m.ID,
		AccountID: m.AccountID,
		Amount:    m.Amount,
		Type:      domain.TransactionType(m.Type),
		CreatedAt: m.CreatedAt,
	}
}

func fromDomainTransaction(t *domain.Transaction) *genmodel.Transactions {
	return &genmodel.Transactions{
		ID:        t.ID,
		AccountID: t.AccountID,
		Amount:    t.Amount,
		Type:      string(t.Type),
		CreatedAt: t.CreatedAt,
	}
}
