package userwallet

import (
	"context"
	"database/sql"
	"fmt"

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
func (r UserWalletRepositoryImpl) GetUserWallet(ctx context.Context, opts ...sqkit.SelectOption) (wallet entity.UserWallet, err error) {
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
func (r UserWalletRepositoryImpl) GetUserWalletByUserXID(ctx context.Context, userXID string) (wallet entity.UserWallet, err error) {
	return r.GetUserWallet(ctx, []sqkit.SelectOption{
		sqkit.Eq{fmt.Sprintf("%s.%s", entity.UserWalletTableName, entity.UserWalletTable.UserXID): userXID},
	}...)
}
