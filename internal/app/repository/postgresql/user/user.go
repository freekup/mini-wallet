package user

import (
	"context"

	"github.com/freekup/mini-wallet/internal/app/entity"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

// @mock
type UserRepository interface {
	GetUser(ctx context.Context, opts ...sqkit.SelectOption) (user entity.User, err error)
	GetActiveUser(ctx context.Context, opts ...sqkit.SelectOption) (user entity.User, err error)
	GetActiveUserByXID(ctx context.Context, xid string) (user entity.User, err error)
}
