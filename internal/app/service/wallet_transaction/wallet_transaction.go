package wallettransaction

import (
	"context"

	"github.com/freekup/mini-wallet/internal/app/entity"
	"github.com/freekup/mini-wallet/pkg/cerror"
)

type (
	WalletTransactionService interface {
		GetWalletTransactions(ctx context.Context, pagination entity.ViewPagination, userXID string) (results []entity.WalletTransaction, pg entity.ViewPagination, cerr cerror.CError)
		AddBalanceWallet(ctx context.Context, arg entity.AddBalanceWalletArg) (walletTransaction entity.WalletTransaction, cerr cerror.CError)
		WithdrawBalance(ctx context.Context, arg entity.WithdrawBalanceArg) (walletTransaction entity.WalletTransaction, cerr cerror.CError)
	}
)
