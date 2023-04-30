package userwallet

import (
	"context"
	"database/sql"
	"errors"

	"github.com/freekup/mini-wallet/internal/app/entity"
	ur "github.com/freekup/mini-wallet/internal/app/repository/postgresql/user"
	uwr "github.com/freekup/mini-wallet/internal/app/repository/postgresql/user_wallet"
	"go.uber.org/dig"
)

type (
	UserWalletServiceImpl struct {
		dig.In
		UserRepo       ur.UserRepository
		UserWalletRepo uwr.UserWalletRepository
	}
)

// @ctor
func NewUserWalletService(impl UserWalletServiceImpl) UserWalletService {
	return &impl
}

// InitializeWallet used to create new wallet for user
func (s *UserWalletServiceImpl) InitializeWallet(ctx context.Context, xid string) (wallet entity.UserWallet, err error) {
	if xid == "" {
		err = errors.New("user xid is empty")
		return
	}

	user, err := s.UserRepo.GetActiveUserByXID(ctx, xid)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if user.ID == 0 {
		err = errors.New("user not found")
		return
	}

	wallet, err = s.UserWalletRepo.GetUserWalletByUserXID(ctx, false, user.XID)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if wallet.ID == "" {
		wallet, err = s.UserWalletRepo.CreateUserWallet(ctx, entity.CreateUserWalletArg{
			UserXID: user.XID,
		})
		if err != nil {
			return
		}
	}

	return
}
