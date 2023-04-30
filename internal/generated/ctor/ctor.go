package ctor

/* DO NOT EDIT. This file generated due to '@ctor' annotation*/

import (
	a "github.com/freekup/mini-wallet/internal/app/infra"
	"github.com/typical-go/typical-go/pkg/typapp"
)

func init() {
	typapp.Provide("", a.NewDatabases)
	typapp.Provide("", a.NewEcho)
}
