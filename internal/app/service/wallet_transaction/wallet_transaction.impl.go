package wallettransaction

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/freekup/mini-wallet/internal/app/entity"
	"github.com/freekup/mini-wallet/internal/app/repository/kafka"
	uwr "github.com/freekup/mini-wallet/internal/app/repository/postgresql/user_wallet"
	wtr "github.com/freekup/mini-wallet/internal/app/repository/postgresql/wallet_transaction"
	"github.com/freekup/mini-wallet/pkg/cerror"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
	"go.uber.org/dig"
)

type (
	WalletTransactionServiceImpl struct {
		dig.In
		WalletTransRepo wtr.WalletTransactionRepository
		UserWalletRepo  uwr.UserWalletRepository
		MessageRepo     kafka.MessageRepository
	}
)

// @ctor
func NewWalletTransactionService(impl WalletTransactionServiceImpl) WalletTransactionService {
	return &impl
}

// GetWalletTransactions used to get list of wallet transactions
func (s *WalletTransactionServiceImpl) GetWalletTransactions(ctx context.Context, pagination entity.ViewPagination, userXID string) (results []entity.WalletTransaction, pg entity.ViewPagination, cerr cerror.CError) {
	var (
		err error
	)

	defer func() {
		if err != nil {
			cerr = cerror.NewSystemError(err.Error())
		}
	}()

	if pagination.Limit < 0 {
		cerr = cerror.NewValidationError("limit=the limit must not be negative")
		return
	} else if pagination.Limit > 20 {
		cerr = cerror.NewValidationError("limit=the maximum limit is 20, the limit must be less or equal to 20")
		return
	}

	if pagination.Offset < 0 {
		cerr = cerror.NewValidationError("offset=the offset must not be negative")
		return
	}
	if userXID == "" {
		cerr = cerror.NewValidationError("xid=user XID is empty")
		return
	}

	opts := make([]sqkit.SelectOption, 0)
	if userXID != "" {
		opts = append(opts, sqkit.Eq{"uw.user_xid": userXID})
	}

	results, pg, err = s.WalletTransRepo.GetWalletTransactions(ctx, pagination, opts...)
	if err != nil {
		return
	}

	return
}

// AddBalanceWallet used to add balance amount
func (s *WalletTransactionServiceImpl) AddBalanceWallet(ctx context.Context, arg entity.AddBalanceWalletArg) (walletTransaction entity.WalletTransaction, cerr cerror.CError) {
	var (
		err error
	)

	defer func() {
		if err != nil {
			cerr = cerror.NewSystemError(err.Error())
		}
	}()

	if arg.ReferenceID == "" {
		cerr = cerror.NewValidationError("reference_id=reference id is empty")
		return
	}
	if arg.Amount <= 0 {
		cerr = cerror.NewValidationError("amount=invalid amount")
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
		cerr = cerror.NewValidationError("user=user wallet not found")
		return
	}
	if !userWallet.IsEnabledBool() {
		cerr = cerror.NewValidationError("user=user wallet is disabled")
		return
	}

	defer func() {
		if err == nil {
			s.MessageRepo.CreatedWalletTransaction(entity.KafkaCreatedWalletTransactionData{
				UserXID: userWallet.UserXID,
				Amount:  arg.Amount,
			})
		}
	}()

	// Update user-wallet current balance
	userWallet.CurrentBalance += arg.Amount
	err = s.UserWalletRepo.UpdateWalletCurrentBalance(ctx, userWallet)
	if err != nil {
		return
	}

	walletTransaction, err = s.WalletTransRepo.CreateWalletTransaction(ctx, entity.WalletTransaction{
		WalletID:    userWallet.ID,
		CreatedBy:   arg.Requestor,
		Status:      entity.Success,
		ReferenceID: arg.ReferenceID,
		Amount:      arg.Amount,
		Description: fmt.Sprintf("Deposit at %s", time.Now().Format(time.RFC3339)),
	})
	if err != nil {
		return
	}

	return
}

// WithdrawBalance used to decreate balance amount
func (s *WalletTransactionServiceImpl) WithdrawBalance(ctx context.Context, arg entity.WithdrawBalanceArg) (walletTransaction entity.WalletTransaction, cerr cerror.CError) {
	var (
		err error
	)

	defer func() {
		if err != nil {
			cerr = cerror.NewSystemError(err.Error())
		}
	}()

	if arg.ReferenceID == "" {
		cerr = cerror.NewValidationError("reference_id=reference id is empty")
		return
	}
	if arg.Amount >= 0 {
		cerr = cerror.NewValidationError("amount=invalid amount")
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
		cerr = cerror.NewValidationError("user=user wallet not found")
		return
	}
	if !userWallet.IsEnabledBool() {
		cerr = cerror.NewValidationError("user=user wallet is disabled")
		return
	}

	if math.Abs(arg.Amount) > userWallet.CurrentBalance {
		cerr = cerror.NewValidationError("amount=amount cannot more than current balance")
		return
	}

	defer func() {
		if err == nil {
			s.MessageRepo.CreatedWalletTransaction(entity.KafkaCreatedWalletTransactionData{
				UserXID: userWallet.UserXID,
				Amount:  arg.Amount,
			})
		}
	}()

	// Update user-wallet current balance
	userWallet.CurrentBalance += arg.Amount
	err = s.UserWalletRepo.UpdateWalletCurrentBalance(ctx, userWallet)
	if err != nil {
		return
	}

	walletTransaction, err = s.WalletTransRepo.CreateWalletTransaction(ctx, entity.WalletTransaction{
		WalletID:    userWallet.ID,
		CreatedBy:   arg.Requestor,
		Status:      entity.Success,
		ReferenceID: arg.ReferenceID,
		Amount:      arg.Amount,
		Description: fmt.Sprintf("Withdraw at %s", time.Now().Format(time.RFC3339)),
	})
	if err != nil {
		return
	}

	return
}
