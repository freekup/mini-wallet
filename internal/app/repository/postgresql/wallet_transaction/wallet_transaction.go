package wallettransaction

import (
	"context"

	"github.com/freekup/mini-wallet/internal/app/entity"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

type (
	// @mock
	WalletTransactionRepository interface {
		GetWalletTransactions(ctx context.Context, pagination entity.ViewPagination, opts ...sqkit.SelectOption) (results []entity.WalletTransaction, pag entity.ViewPagination, err error)
		CreateWalletTransaction(ctx context.Context, wTransaction entity.WalletTransaction) (result entity.WalletTransaction, err error)
	}
)
