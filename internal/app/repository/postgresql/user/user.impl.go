package user

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

type UserRepositoryImpl struct {
	dig.In
	*sql.DB `name:"pg"`
}

// @ctor
func NewUserRepository(impl UserRepositoryImpl) UserRepository {
	return &impl
}

// GetUser used to get user data, pass opts to add new condition
func (r *UserRepositoryImpl) GetUser(ctx context.Context, opts ...sqkit.SelectOption) (user entity.User, err error) {
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
		fmt.Sprintf("%s.%s", entity.UserTableName, entity.UserTable.ID),
		fmt.Sprintf("%s.%s", entity.UserTableName, entity.UserTable.Name),
		fmt.Sprintf("%s.%s", entity.UserTableName, entity.UserTable.XID),
	).From(entity.UserTableName).PlaceholderFormat(sq.Dollar)

	for _, opt := range opts {
		queryBuilder = opt.CompileSelect(queryBuilder)
	}

	err = queryBuilder.RunWith(txn).QueryRowContext(ctx).Scan(
		&user.ID,
		&user.Name,
		&user.XID,
	)
	if err != nil {
		return
	}

	return
}

// GetActiveUser used to get only user with is_active status is 1
func (r *UserRepositoryImpl) GetActiveUser(ctx context.Context, opts ...sqkit.SelectOption) (user entity.User, err error) {
	opts = append(opts, sqkit.Eq{fmt.Sprintf("%s.%s", entity.UserTableName, entity.UserTable.IsActive): 1})
	return r.GetUser(ctx, opts...)
}

// GetActiveUser used to get only user with is_active status is 1 and filtered by given XID
func (r *UserRepositoryImpl) GetActiveUserByXID(ctx context.Context, xid string) (user entity.User, err error) {
	return r.GetActiveUser(ctx, []sqkit.SelectOption{
		sqkit.Eq{fmt.Sprintf("%s.%s", entity.UserTableName, entity.UserTable.XID): xid},
	}...)
}
