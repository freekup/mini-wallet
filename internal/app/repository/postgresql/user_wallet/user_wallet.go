package userwallet

import (
	"context"

	"github.com/freekup/mini-wallet/internal/app/entity"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

type (
	// @mock
	UserWalletRepository interface {
		GetUserWallet(ctx context.Context, opts ...sqkit.SelectOption) (wallet entity.UserWallet, err error)
		GetUserWalletByUserXID(ctx context.Context, userXID string) (wallet entity.UserWallet, err error)
	}
)
