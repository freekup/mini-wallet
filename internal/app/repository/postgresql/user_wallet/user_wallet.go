package userwallet

import (
	"context"

	"github.com/freekup/mini-wallet/internal/app/entity"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

type (
	// @mock
	UserWalletRepository interface {
		GetUserWallet(ctx context.Context, isLock bool, opts ...sqkit.SelectOption) (wallet entity.UserWallet, err error)
		GetUserWalletByUserXID(ctx context.Context, isLock bool, userXID string) (wallet entity.UserWallet, err error)
		CreateUserWallet(ctx context.Context, arg entity.CreateUserWalletArg) (wallet entity.UserWallet, err error)
		ChangeEnableStatusWallet(ctx context.Context, wallet entity.UserWallet) (err error)
		UpdateWalletCurrentBalance(ctx context.Context, wallet entity.UserWallet) (err error)
	}
)
