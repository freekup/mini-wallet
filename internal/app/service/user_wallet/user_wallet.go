package userwallet

import (
	"context"

	"github.com/freekup/mini-wallet/internal/app/entity"
	"github.com/freekup/mini-wallet/pkg/cerror"
)

type (
	UserWalletService interface {
		InitializeWallet(ctx context.Context, xid string) (wallet entity.UserWallet, cerr cerror.CError)
		EnableWallet(ctx context.Context, userXID string) (wallet entity.UserWallet, cerr cerror.CError)
		DisableWallet(ctx context.Context, isDisable bool, userXID string) (wallet entity.UserWallet, cerr cerror.CError)
		GetUserWalletByUserXID(ctx context.Context, userXID string) (wallet entity.UserWallet, cerr cerror.CError)
		RefreshUserWalletCache(ctx context.Context, userXID string) (err error)
	}
)
