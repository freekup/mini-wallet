package userwallet

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/freekup/mini-wallet/internal/app/entity"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
	"go.uber.org/dig"
)

type (
	UserWalletRepositoryImpl struct {
		dig.In
		*sql.DB `name:"pg"`
	}
)

// @ctor
func NewUserWalletRepository(impl UserWalletRepositoryImpl) UserWalletRepository {
	return &impl
}

// GetUserWallet used to get user wallet data with given condition
// isLock used to identify is it LOCK query or not
func (r UserWalletRepositoryImpl) GetUserWallet(ctx context.Context, isLock bool, opts ...sqkit.SelectOption) (wallet entity.UserWallet, err error) {
	// Initialize transaction session from context
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	defer func() {
		if err != nil && txn != nil {
			// Will automatic rollback if any error
			txn.AppendError(err)
		}
	}()

	queryBuilder := sq.Select(
		fmt.Sprintf("%s.%s", entity.UserWalletTableName, entity.UserWalletTable.ID),
		fmt.Sprintf("%s.%s", entity.UserWalletTableName, entity.UserWalletTable.UserXID),
		fmt.Sprintf("%s.%s", entity.UserWalletTableName, entity.UserWalletTable.CurrentBalance),
		fmt.Sprintf("%s.%s", entity.UserWalletTableName, entity.UserWalletTable.IsEnabled),
		fmt.Sprintf("TO_CHAR(%s.%s, 'YYYY-MM-DD\"T\"HH:mm:ssTZH')", entity.UserWalletTableName, entity.UserWalletTable.EnabledAt),
		fmt.Sprintf("%s.%s", entity.UserWalletTableName, entity.UserWalletTable.DeletedBy),
		fmt.Sprintf("TO_CHAR(%s.%s, 'YYYY-MM-DD\"T\"HH:mm:ssTZH')", entity.UserWalletTableName, entity.UserWalletTable.DeletedAt),
	).From(entity.UserWalletTableName).PlaceholderFormat(sq.Dollar)

	for _, opt := range opts {
		queryBuilder = opt.CompileSelect(queryBuilder)
	}

	// This query used to locking data from other connection session
	// usually used to prevent race condition
	if isLock {
		queryBuilder.Suffix("FOR UPDATE")
	}

	err = queryBuilder.RunWith(txn).QueryRowContext(ctx).Scan(
		&wallet.ID,
		&wallet.UserXID,
		&wallet.CurrentBalance,
		&wallet.IsEnabled,
		&wallet.EnabledAt,
		&wallet.DeletedBy,
		&wallet.DeletedAt,
	)
	if err != nil {
		return
	}

	return
}

// GetUserWalletByUserID used to get user wallet with user_id filter
func (r UserWalletRepositoryImpl) GetUserWalletByUserXID(ctx context.Context, isLock bool, userXID string) (wallet entity.UserWallet, err error) {
	return r.GetUserWallet(ctx, isLock, []sqkit.SelectOption{
		sqkit.Eq{fmt.Sprintf("%s.%s", entity.UserWalletTableName, entity.UserWalletTable.UserXID): userXID},
	}...)
}

// CreateUserWallet used to create new user wallet data
func (r UserWalletRepositoryImpl) CreateUserWallet(ctx context.Context, arg entity.CreateUserWalletArg) (wallet entity.UserWallet, err error) {
	// Initialize transaction session from context
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	defer func() {
		if err != nil && txn != nil {
			// Will automatic rollback if any error
			txn.AppendError(err)
		}
	}()

	builder := sq.Insert(entity.UserWalletTableName).Columns(
		entity.UserWalletTable.ID,
		entity.UserWalletTable.UserXID,
		entity.UserWalletTable.CurrentBalance,
		entity.UserWalletTable.IsEnabled,
	).Values(
		sq.Expr("gen_random_uuid()"),
		arg.UserXID,
		0,
		0,
	).Suffix(fmt.Sprintf("RETURNING %s, %s, %s, %s",
		entity.UserWalletTable.ID,
		entity.UserWalletTable.UserXID,
		entity.UserWalletTable.CurrentBalance,
		entity.UserWalletTable.IsEnabled,
	)).PlaceholderFormat(sq.Dollar)

	rows, err := builder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&wallet.ID,
			&wallet.UserXID,
			&wallet.CurrentBalance,
			&wallet.IsEnabled,
		); err != nil {
			return
		}
	}

	return
}

// ChaneEnableStatusWallet used to change status from user wallet
func (r *UserWalletRepositoryImpl) ChangeEnableStatusWallet(ctx context.Context, wallet entity.UserWallet) (err error) {
	// Initialize transaction session from context
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	defer func() {
		if err != nil && txn != nil {
			// Will automatic rollback if any error
			txn.AppendError(err)
		}
	}()

	builder := sq.Update(entity.UserWalletTableName).
		Set(entity.UserWalletTable.IsEnabled, wallet.IsEnabled).
		Set(entity.UserWalletTable.EnabledAt, wallet.EnabledAt).
		Set(entity.UserWalletTable.ModifiedAt, time.Now().Format(time.RFC3339)).
		Set(entity.UserWalletTable.DeletedAt, wallet.DeletedAt).
		Set(entity.UserWalletTable.DeletedBy, wallet.DeletedBy).
		Where(sq.Eq{entity.UserWalletTable.ID: wallet.ID}).PlaceholderFormat(sq.Dollar)

	_, err = builder.RunWith(txn).ExecContext(ctx)
	if err != nil {
		return
	}

	return
}

// UpdateWalletCurrentBalance used to change current balance from user wallet
func (r UserWalletRepositoryImpl) UpdateWalletCurrentBalance(ctx context.Context, wallet entity.UserWallet) (err error) {
	// Initialize transaction session from context
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	defer func() {
		if err != nil && txn != nil {
			// Will automatic rollback if any error
			txn.AppendError(err)
		}
	}()

	builder := sq.Update(entity.UserWalletTableName).
		Set(entity.UserWalletTable.CurrentBalance, wallet.CurrentBalance).
		Set(entity.UserWalletTable.ModifiedAt, time.Now().Format(time.RFC3339)).
		Where(sq.Eq{entity.UserWalletTable.ID: wallet.ID}).PlaceholderFormat(sq.Dollar)

	_, err = builder.RunWith(txn).ExecContext(ctx)
	if err != nil {
		return
	}

	return
}
