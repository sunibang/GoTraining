package service_test

import (
	"context"
	"testing"

	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/repository/mocks"
	"github.com/romangurevitch/go-training/internal/bank/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBankService_Deposit(t *testing.T) {
	// QUEST 6: Participants should implement this table-driven test.

	type args struct {
		accountOwner  string
		depositAmount int64
	}

	tests := []struct {
		name            string
		args            args
		wantErr         bool
		expectedBalance int64
	}{
		{
			name: "Successful deposit",
			args: args{
				accountOwner:  "John Doe",
				depositAmount: 10000, // 100.00
			},
			wantErr:         false,
			expectedBalance: 10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			bs := service.NewBankService(repo)
			ctx := context.Background()

			acc := &domain.Account{
				ID:      "ACC-1",
				Owner:   tt.args.accountOwner,
				Balance: 0,
				Status:  domain.StatusOpen,
			}

			// Expect GetAccount for verification after deposit
			repo.EXPECT().GetAccount(ctx, "ACC-1").Return(acc, nil)
			repo.EXPECT().SaveAccount(ctx, mock.Anything).Return(nil)
			repo.EXPECT().SaveTransaction(ctx, mock.Anything).Return(nil)

			// Perform deposit
			err := bs.Deposit(ctx, acc.ID, tt.args.depositAmount)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBalance, acc.Balance)
			}
		})
	}
}
