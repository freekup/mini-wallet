package ctor

/* DO NOT EDIT. This file generated due to '@ctor' annotation*/

import (
	a "github.com/freekup/mini-wallet/internal/app/infra"
	b "github.com/freekup/mini-wallet/internal/app/repository/postgresql/user"
	c "github.com/freekup/mini-wallet/internal/app/repository/postgresql/user_wallet"
	d "github.com/freekup/mini-wallet/internal/app/service/user_wallet"
	"github.com/typical-go/typical-go/pkg/typapp"
)

func init() {
	typapp.Provide("", a.NewDatabases)
	typapp.Provide("", a.NewEcho)
	typapp.Provide("", b.NewUserRepository)
	typapp.Provide("", c.NewUserWalletRepository)
	typapp.Provide("", d.NewUserWalletService)
}
