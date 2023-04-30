package rest

import (
	"github.com/freekup/mini-wallet/internal/app/middleware"
	uws "github.com/freekup/mini-wallet/internal/app/service/user_wallet"
	"github.com/freekup/mini-wallet/pkg/tokenizer"
	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"go.uber.org/dig"
)

type (
	UserWalletController struct {
		dig.In
		UserWalletService uws.UserWalletService
	}
)

var _ echokit.Router = (*UserWalletController)(nil)

func (c *UserWalletController) SetRoute(e echokit.Server) {
	e.POST("/init", c.InitializeWallet)

	wallet := e.Group("/wallet")
	wallet.POST("", c.EnableWallet, middleware.JWTAuth)
}

// InitializeWallet used to handle InitializeWallet Rest
// also to generate token to access other API
func (c *UserWalletController) InitializeWallet(ec echo.Context) (err error) {
	ctx := ec.Request().Context()
	param := struct {
		UserXID string `form:"customer_xid"`
	}{}

	if err = ec.Bind(&param); err != nil {
		return
	}

	userWallet, err := c.UserWalletService.InitializeWallet(ctx, param.UserXID)
	if err != nil {
		return
	}

	token, err := tokenizer.GenerateJWTToken(ctx, userWallet.UserXID)
	if err != nil {
		return
	}

	return ec.JSON(200, map[string]interface{}{
		"token": token,
	})
}

// EnableWallet used to enable wallet status
func (c *UserWalletController) EnableWallet(ec echo.Context) (err error) {
	var (
		ctx = ec.Request().Context()
		xid = tokenizer.FetchEchoTokenXID(ec)
	)

	wallet, err := c.UserWalletService.EnableWallet(ctx, xid)
	if err != nil {
		return
	}

	return ec.JSON(200, map[string]interface{}{
		"id":         wallet.ID,
		"owned_by":   wallet.UserXID,
		"status":     wallet.IsEnabledString(),
		"enabled_at": wallet.EnabledAt,
		"balance":    wallet.CurrentBalance,
	})
}
