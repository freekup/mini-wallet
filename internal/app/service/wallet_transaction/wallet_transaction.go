package wallettransaction

import (
	"context"

	"github.com/freekup/mini-wallet/internal/app/entity"
)

type (
	WalletTransactionService interface {
		AddBalanceWallet(ctx context.Context, arg entity.AddBalanceWalletArg) (walletTransaction entity.WalletTransaction, err error)
		WithdrawBalance(ctx context.Context, arg entity.WithdrawBalanceArg) (walletTransaction entity.WalletTransaction, err error)
	}
)
