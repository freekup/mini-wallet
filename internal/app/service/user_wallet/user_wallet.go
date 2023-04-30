package userwallet

import (
	"context"

	"github.com/freekup/mini-wallet/internal/app/entity"
)

type (
	UserWalletService interface {
		InitializeWallet(ctx context.Context, xid string) (wallet entity.UserWallet, err error)
	}
)
