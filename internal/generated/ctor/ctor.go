package ctor

/* DO NOT EDIT. This file generated due to '@ctor' annotation*/

import (
	a "github.com/freekup/mini-wallet/internal/app/infra"
	b "github.com/freekup/mini-wallet/internal/app/repository/kafka"
	c "github.com/freekup/mini-wallet/internal/app/repository/postgresql/user"
	d "github.com/freekup/mini-wallet/internal/app/repository/postgresql/user_wallet"
	e "github.com/freekup/mini-wallet/internal/app/repository/postgresql/wallet_transaction"
	f "github.com/freekup/mini-wallet/internal/app/service/user_wallet"
	g "github.com/freekup/mini-wallet/internal/app/service/wallet_transaction"
	"github.com/typical-go/typical-go/pkg/typapp"
)

func init() {
	typapp.Provide("", a.NewCacheStore)
	typapp.Provide("", a.NewDatabases)
	typapp.Provide("", a.NewEcho)
	typapp.Provide("", a.NewProducer)
	typapp.Provide("", a.NewConsumer)
	typapp.Provide("", b.NewMessageRepository)
	typapp.Provide("", c.NewUserRepository)
	typapp.Provide("", d.NewUserWalletRepository)
	typapp.Provide("", e.NewWalletTransactionRepository)
	typapp.Provide("", f.NewUserWalletService)
	typapp.Provide("", g.NewWalletTransactionService)
}
