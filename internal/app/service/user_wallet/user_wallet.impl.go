package userwallet

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/freekup/mini-wallet/internal/app/entity"
	ur "github.com/freekup/mini-wallet/internal/app/repository/postgresql/user"
	uwr "github.com/freekup/mini-wallet/internal/app/repository/postgresql/user_wallet"
	"github.com/freekup/mini-wallet/pkg/cerror"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
	"go.uber.org/dig"
)

type (
	UserWalletServiceImpl struct {
		dig.In
		UserRepo       ur.UserRepository
		UserWalletRepo uwr.UserWalletRepository
		Redis          *cachekit.Store
	}
)

// @ctor
func NewUserWalletService(impl UserWalletServiceImpl) UserWalletService {
	return &impl
}

// InitializeWallet used to create new wallet for user
func (s *UserWalletServiceImpl) InitializeWallet(ctx context.Context, xid string) (wallet entity.UserWallet, cerr cerror.CError) {
	var (
		err error
	)

	defer func() {
		if err != nil && cerr == nil {
			cerr = cerror.NewSystemError(err.Error())
		}
	}()

	if xid == "" {
		cerr = cerror.NewValidationError("user_xid=user xid is empty")
		return
	}

	user, err := s.UserRepo.GetActiveUserByXID(ctx, xid)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if user.ID == 0 {
		cerr = cerror.NewValidationError("user=user not found")
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

// EnableWallet used to change wallet status from disable to enable
func (s *UserWalletServiceImpl) EnableWallet(ctx context.Context, userXID string) (wallet entity.UserWallet, cerr cerror.CError) {
	var (
		err error
	)

	defer func() {
		if err != nil && cerr == nil {
			cerr = cerror.NewSystemError(err.Error())
		}
	}()

	if userXID == "" {
		cerr = cerror.NewValidationError("xid=user xid is empty")
		return
	}

	wallet, err = s.UserWalletRepo.GetUserWalletByUserXID(ctx, false, userXID)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if wallet.ID == "" {
		cerr = cerror.NewValidationError("wallet=wallet not found")
		return
	}
	if wallet.IsEnabledBool() {
		cerr = cerror.NewValidationError("wallet=wallet already enabled")
		return
	}

	currTime := time.Now().Format(time.RFC3339)

	wallet.IsEnabled = 1
	wallet.EnabledAt = &currTime
	wallet.DeletedAt = nil
	wallet.DeletedBy = nil

	err = s.UserWalletRepo.ChangeEnableStatusWallet(ctx, wallet)
	if err != nil {
		return
	}

	return
}

// DisableWallet used to disable user wallet status
func (s *UserWalletServiceImpl) DisableWallet(ctx context.Context, isDisable bool, userXID string) (wallet entity.UserWallet, cerr cerror.CError) {
	var (
		err error
	)

	defer func() {
		if err != nil && cerr == nil {
			cerr = cerror.NewSystemError(err.Error())
		}
	}()

	if userXID == "" {
		cerr = cerror.NewValidationError("xid=user xid is empty")
		return
	}

	wallet, err = s.UserWalletRepo.GetUserWalletByUserXID(ctx, true, userXID)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if wallet.ID == "" {
		cerr = cerror.NewValidationError("wallet=wallet not found")
		return
	}
	if !isDisable {
		cerr = cerror.NewValidationError("is_disable=is_disable is false")
		return
	}
	if isDisable && !wallet.IsEnabledBool() {
		cerr = cerror.NewValidationError("wallet=wallet already disabled")
		return
	}

	currTime := time.Now().Format(time.RFC3339)

	wallet.IsEnabled = 0
	wallet.EnabledAt = nil
	wallet.DeletedBy = &userXID
	wallet.DeletedAt = &currTime

	err = s.UserWalletRepo.ChangeEnableStatusWallet(ctx, wallet)
	if err != nil {
		return
	}

	s.RemoveWalletCache(ctx, userXID)

	return
}

// GetUserWalletByUserXID used to get user wallet
func (s *UserWalletServiceImpl) GetUserWalletByUserXID(ctx context.Context, userXID string) (wallet entity.UserWallet, cerr cerror.CError) {
	var (
		err error
	)

	defer func() {
		if err != nil && cerr == nil {
			cerr = cerror.NewSystemError(err.Error())
		}
	}()

	if userXID == "" {
		cerr = cerror.NewValidationError("xid=user xid is empty")
		return
	}

	// Fetch from cache
	cacheWallet, err := s.Redis.HGet(ctx, entity.RedisUserBalanceKey, userXID).Result()
	if err != nil && err.Error() != entity.ErrRedisNoRows {
		return
	}

	if cacheWallet == "" {
		wallet, err = s.UserWalletRepo.GetUserWalletByUserXID(ctx, false, userXID)
		if err != nil && err != sql.ErrNoRows {
			return
		}
		if wallet.ID == "" {
			cerr = cerror.NewValidationError("wallet=wallet not found")
			return
		}
		if !wallet.IsEnabledBool() {
			cerr = cerror.NewValidationError("wallet=wallet is disabled")
			return
		}

		// Store value to redis cache
		var bWallet []byte
		bWallet, err = json.Marshal(wallet)
		if err != nil {
			return
		}

		cacheWallet = string(bWallet)
		s.Redis.HSet(ctx, entity.RedisUserBalanceKey, map[string]interface{}{
			userXID: cacheWallet,
		})
	}

	err = json.Unmarshal([]byte(cacheWallet), &wallet)
	if err != nil {
		return
	}

	go s.RefreshUserWalletCache(context.Background(), userXID)

	return
}

// RefreshUserWalletCache used to refresh user wallet cache
func (s *UserWalletServiceImpl) RefreshUserWalletCache(ctx context.Context, userXID string) (err error) {
	if userXID == "" {
		err = errors.New("user xid is empty")
		return
	}

	wallet, err := s.UserWalletRepo.GetUserWalletByUserXID(ctx, false, userXID)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if wallet.ID == "" {
		err = errors.New("wallet not found")
		return
	}
	if !wallet.IsEnabledBool() {
		err = errors.New("wallet is disabled")
		return
	}

	var bWallet []byte
	bWallet, err = json.Marshal(wallet)
	if err != nil {
		return
	}

	cacheWallet := string(bWallet)
	s.Redis.HSet(ctx, entity.RedisUserBalanceKey, map[string]interface{}{
		userXID: cacheWallet,
	})

	return
}

// RemoveWalletCache used to delete existing cache
func (s *UserWalletServiceImpl) RemoveWalletCache(ctx context.Context, userXID string) {
	s.Redis.HDel(ctx, entity.RedisUserBalanceKey, userXID)
}
