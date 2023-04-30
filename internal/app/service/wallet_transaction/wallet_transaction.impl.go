package wallettransaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/freekup/mini-wallet/internal/app/entity"
	uwr "github.com/freekup/mini-wallet/internal/app/repository/postgresql/user_wallet"
	wtr "github.com/freekup/mini-wallet/internal/app/repository/postgresql/wallet_transaction"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
	"go.uber.org/dig"
)

type (
	WalletTransactionServiceImpl struct {
		dig.In
		WalletTransRepo wtr.WalletTransactionRepository
		UserWalletRepo  uwr.UserWalletRepository
	}
)

// @ctor
func NewWalletTransactionService(impl WalletTransactionServiceImpl) WalletTransactionService {
	return &impl
}

// GetWalletTransactions used to get list of wallet transactions
func (s *WalletTransactionServiceImpl) GetWalletTransactions(ctx context.Context, pagination entity.ViewPagination, userXID string) (results []entity.WalletTransaction, pg entity.ViewPagination, err error) {
	if pagination.Limit < 0 {
		err = errors.New("the limit must not be negative")
		return
	} else if pagination.Limit > 20 {
		err = errors.New("the maximum limit is 20, the limit must be less or equal to 20")
		return
	}

	if pagination.Offset < 0 {
		err = errors.New("the offset must not be negative")
		return
	}
	if userXID == "" {
		err = errors.New("user XID is empty")
		return
	}

	opts := make([]sqkit.SelectOption, 0)
	if userXID != "" {
		opts = append(opts, sqkit.Eq{"uw.user_xid": userXID})
	}

	return s.WalletTransRepo.GetWalletTransactions(ctx, pagination, opts...)
}

// AddBalanceWallet used to add balance amount
func (s *WalletTransactionServiceImpl) AddBalanceWallet(ctx context.Context, arg entity.AddBalanceWalletArg) (walletTransaction entity.WalletTransaction, err error) {
	if arg.ReferenceID == "" {
		err = errors.New("reference id is empty")
		return
	}
	if arg.Amount <= 0 {
		err = errors.New("invalid amount")
		return
	}

	// Begin Transaction
	tx := dbtxn.Begin(&ctx)
	defer func() {
		if err != nil {
			tx.AppendError(err)
		} else {
			tx.Commit()
		}
	}()

	// Use lock query to prevent RACE CONDITION issue
	userWallet, err := s.UserWalletRepo.GetUserWalletByUserXID(ctx, true, arg.Requestor)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if userWallet.ID == "" {
		err = errors.New("user wallet not found")
		return
	}
	if !userWallet.IsEnabledBool() {
		err = errors.New("user wallet is disabled")
		return
	}

	// Update user-wallet current balance
	userWallet.CurrentBalance += arg.Amount
	err = s.UserWalletRepo.UpdateWalletCurrentBalance(ctx, userWallet)
	if err != nil {
		return
	}

	return s.WalletTransRepo.CreateWalletTransaction(ctx, entity.WalletTransaction{
		WalletID:    userWallet.ID,
		CreatedBy:   arg.Requestor,
		Status:      entity.Success,
		ReferenceID: arg.ReferenceID,
		Amount:      arg.Amount,
		Description: fmt.Sprintf("Deposit at %s", time.Now().Format(time.RFC3339)),
	})
}

// WithdrawBalance used to decreate balance amount
func (s *WalletTransactionServiceImpl) WithdrawBalance(ctx context.Context, arg entity.WithdrawBalanceArg) (walletTransaction entity.WalletTransaction, err error) {
	if arg.ReferenceID == "" {
		err = errors.New("reference id not found")
		return
	}
	if arg.Amount >= 0 {
		err = errors.New("invalid amount")
		return
	}

	// Begin Transaction
	tx := dbtxn.Begin(&ctx)
	defer func() {
		if err != nil {
			tx.AppendError(err)
		} else {
			tx.Commit()
		}
	}()

	// Use lock query to prevent RACE CONDITION issue
	userWallet, err := s.UserWalletRepo.GetUserWalletByUserXID(ctx, true, arg.Requestor)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if userWallet.ID == "" {
		err = errors.New("user wallet not found")
		return
	}
	if !userWallet.IsEnabledBool() {
		err = errors.New("user wallet is disabled")
		return
	}

	if math.Abs(arg.Amount) > userWallet.CurrentBalance {
		err = errors.New("amount cannot more than current balance")
		return
	}

	// Update user-wallet current balance
	userWallet.CurrentBalance += arg.Amount
	err = s.UserWalletRepo.UpdateWalletCurrentBalance(ctx, userWallet)
	if err != nil {
		return
	}

	return s.WalletTransRepo.CreateWalletTransaction(ctx, entity.WalletTransaction{
		WalletID:    userWallet.ID,
		CreatedBy:   arg.Requestor,
		Status:      entity.Success,
		ReferenceID: arg.ReferenceID,
		Amount:      arg.Amount,
		Description: fmt.Sprintf("Withdraw at %s", time.Now().Format(time.RFC3339)),
	})
}
