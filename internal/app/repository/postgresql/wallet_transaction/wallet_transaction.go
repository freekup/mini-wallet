package wallettransaction

import (
	"context"

	"github.com/freekup/mini-wallet/internal/app/entity"
)

type (
	// @mock
	WalletTransactionRepository interface {
		CreateWalletTransaction(ctx context.Context, wTransaction entity.WalletTransaction) (result entity.WalletTransaction, err error)
	}
)
