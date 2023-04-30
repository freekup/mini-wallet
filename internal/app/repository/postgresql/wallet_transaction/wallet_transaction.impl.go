package wallettransaction

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/freekup/mini-wallet/internal/app/entity"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"go.uber.org/dig"
)

type (
	WalletTransactionRepositoryImpl struct {
		dig.In
		*sql.DB `name:"pg"`
	}
)

// @ctor
func NewWalletTransactionRepository(impl WalletTransactionRepositoryImpl) WalletTransactionRepository {
	return &impl
}

// CreateWalletTransaction used to create wallet transaction
func (r *WalletTransactionRepositoryImpl) CreateWalletTransaction(ctx context.Context, wTransaction entity.WalletTransaction) (result entity.WalletTransaction, err error) {
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

	builder := sq.Insert(entity.WalletTransactionTableName).Columns(
		entity.WalletTransactionTable.ID,
		entity.WalletTransactionTable.WalletID,
		entity.WalletTransactionTable.Status,
		entity.WalletTransactionTable.ReferenceID,
		entity.WalletTransactionTable.Amount,
		entity.WalletTransactionTable.Description,
		entity.WalletTransactionTable.CreatedBy,
	).Values(
		sq.Expr("gen_random_uuid()"),
		wTransaction.WalletID,
		wTransaction.Status,
		wTransaction.ReferenceID,
		wTransaction.Amount,
		wTransaction.Description,
		wTransaction.CreatedBy,
	).Suffix(fmt.Sprintf("RETURNING %s, %s, %s, %s, %s, %s, %s, TO_CHAR(%s, 'YYYY-MM-DD\"T\"HH:mm:ssTZH')",
		entity.WalletTransactionTable.ID,
		entity.WalletTransactionTable.WalletID,
		entity.WalletTransactionTable.Status,
		entity.WalletTransactionTable.ReferenceID,
		entity.WalletTransactionTable.Amount,
		entity.WalletTransactionTable.Description,
		entity.WalletTransactionTable.CreatedBy,
		entity.WalletTransactionTable.CreatedAt,
	)).PlaceholderFormat(sq.Dollar)

	rows, err := builder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&result.ID,
			&result.WalletID,
			&result.Status,
			&result.ReferenceID,
			&result.Amount,
			&result.Description,
			&result.CreatedBy,
			&result.CreatedAt,
		); err != nil {
			return
		}
	}

	return
}
