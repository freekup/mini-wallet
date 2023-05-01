package wallettransaction

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
	WalletTransactionRepositoryImpl struct {
		dig.In
		*sql.DB `name:"pg"`
	}
)

// @ctor
func NewWalletTransactionRepository(impl WalletTransactionRepositoryImpl) WalletTransactionRepository {
	return &impl
}

// GetWalletTransactions used to get all wallet transaction
func (r *WalletTransactionRepositoryImpl) GetWalletTransactions(ctx context.Context, pagination entity.ViewPagination, opts ...sqkit.SelectOption) (results []entity.WalletTransaction, pag entity.ViewPagination, err error) {
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
		"COUNT(*) OVER() AS total",
		fmt.Sprintf("%s.%s", entity.WalletTransactionTableName, entity.WalletTransactionTable.ID),
		fmt.Sprintf("%s.%s", entity.WalletTransactionTableName, entity.WalletTransactionTable.WalletID),
		fmt.Sprintf("%s.%s", entity.WalletTransactionTableName, entity.WalletTransactionTable.ReferenceID),
		fmt.Sprintf("%s.%s", entity.WalletTransactionTableName, entity.WalletTransactionTable.Amount),
		fmt.Sprintf("%s.%s", entity.WalletTransactionTableName, entity.WalletTransactionTable.Description),
		fmt.Sprintf("%s.%s", entity.WalletTransactionTableName, entity.WalletTransactionTable.CreatedBy),
		fmt.Sprintf("TO_CHAR(%s.%s, 'YYYY-MM-DD\"T\"HH:mm:ssTZH')", entity.WalletTransactionTableName, entity.WalletTransactionTable.CreatedAt),
	).From(entity.WalletTransactionTableName).PlaceholderFormat(sq.Dollar).
		// Use subquery to improve query process
		Join(fmt.Sprintf(`
			(
				SELECT
					id,
					user_xid
				FROM public.user_wallets
			) uw ON uw.id = %s.%s
		`, entity.WalletTransactionTableName, entity.WalletTransactionTable.WalletID))

	for _, opt := range opts {
		queryBuilder = opt.CompileSelect(queryBuilder)
	}

	if pagination.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(pagination.Limit))
	}
	if pagination.Offset > 0 {
		queryBuilder = queryBuilder.Offset(uint64(pagination.Offset))
	}

	rows, err := queryBuilder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		result := entity.WalletTransaction{}
		if err = rows.Scan(
			&pag.Total,
			&result.ID,
			&result.WalletID,
			&result.ReferenceID,
			&result.Amount,
			&result.Description,
			&result.CreatedBy,
			&result.CreatedAt,
		); err != nil {
			return
		}

		results = append(results, result)
	}

	pag.Limit = pagination.Limit
	pag.Offset = pagination.Offset

	return
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
