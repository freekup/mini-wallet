package wallettransaction

import (
	"context"

	"github.com/freekup/mini-wallet/internal/app/entity"
)

type (
	WalletTransactionService interface {
		GetWalletTransactions(ctx context.Context, pagination entity.ViewPagination, userXID string) (results []entity.WalletTransaction, pg entity.ViewPagination, err error)
		AddBalanceWallet(ctx context.Context, arg entity.AddBalanceWalletArg) (walletTransaction entity.WalletTransaction, err error)
		WithdrawBalance(ctx context.Context, arg entity.WithdrawBalanceArg) (walletTransaction entity.WalletTransaction, err error)
	}
)
