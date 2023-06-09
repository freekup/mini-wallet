package rest

import (
	"net/http"

	"github.com/freekup/mini-wallet/internal/app/middleware"
	uws "github.com/freekup/mini-wallet/internal/app/service/user_wallet"
	wts "github.com/freekup/mini-wallet/internal/app/service/wallet_transaction"
	"github.com/freekup/mini-wallet/pkg/cerror"
	"github.com/freekup/mini-wallet/pkg/jsend"
	"github.com/freekup/mini-wallet/pkg/tokenizer"
	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"go.uber.org/dig"
)

type (
	UserWalletController struct {
		dig.In
		UserWalletService  uws.UserWalletService
		WalletTransService wts.WalletTransactionService
	}
)

var _ echokit.Router = (*UserWalletController)(nil)

func (c *UserWalletController) SetRoute(e echokit.Server) {
	v1 := e.Group("/api/v1")
	v1.POST("/init", c.InitializeWallet)

	wallet := v1.Group("/wallet", middleware.JWTAuth)
	wallet.POST("", c.EnableWallet)
	wallet.PATCH("", c.DisableWallet)
	wallet.GET("", c.ViewMyWallet)

	wallet.GET("/transactions", c.GetWalletTransactions)
	wallet.POST("/deposits", c.WalletDeposit)
	wallet.POST("/withdrawals", c.WalletWithdraw)
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

	userWallet, cerr := c.UserWalletService.InitializeWallet(ctx, param.UserXID)
	if cerr != nil {
		return ec.JSON(cerr.GetStatusCode(), jsend.GenerateResponseError(cerr))
	}

	token, err := tokenizer.GenerateJWTToken(ctx, userWallet.UserXID)
	if err != nil {
		cerr = cerror.NewSystemError(err.Error())
		return ec.JSON(cerr.GetStatusCode(), jsend.GenerateResponseError(cerr))
	}

	return ec.JSON(http.StatusOK, jsend.GenerateResponseSuccess(map[string]interface{}{
		"token": token,
	}))
}

// EnableWallet used to enable wallet status
func (c *UserWalletController) EnableWallet(ec echo.Context) (err error) {
	var (
		ctx = ec.Request().Context()
		xid = tokenizer.FetchEchoTokenXID(ec)
	)

	wallet, cerr := c.UserWalletService.EnableWallet(ctx, xid)
	if cerr != nil {
		return ec.JSON(cerr.GetStatusCode(), jsend.GenerateResponseError(cerr))
	}

	return ec.JSON(http.StatusOK, jsend.GenerateResponseSuccess(map[string]interface{}{
		"wallet": map[string]interface{}{
			"id":         wallet.ID,
			"owned_by":   wallet.UserXID,
			"status":     wallet.IsEnabledString(),
			"enabled_at": wallet.EnabledAt,
			"balance":    wallet.CurrentBalance,
		},
	}))
}

// DisableWallet used to disable wallet status
func (c *UserWalletController) DisableWallet(ec echo.Context) (err error) {
	var (
		ctx        = ec.Request().Context()
		xid        = tokenizer.FetchEchoTokenXID(ec)
		isDisabled = false
	)

	if ec.FormValue("is_disabled") == "true" {
		isDisabled = true
	}

	wallet, cerr := c.UserWalletService.DisableWallet(ctx, isDisabled, xid)
	if cerr != nil {
		return ec.JSON(cerr.GetStatusCode(), jsend.GenerateResponseError(cerr))
	}

	return ec.JSON(http.StatusOK, jsend.GenerateResponseSuccess(map[string]interface{}{
		"wallet": map[string]interface{}{
			"id":          wallet.ID,
			"owned_by":    wallet.UserXID,
			"status":      wallet.IsEnabledString(),
			"disabled_at": wallet.DeletedAt,
			"balance":     wallet.CurrentBalance,
		},
	}))
}

// ViewMyWallet used to get wallet from auth
func (c *UserWalletController) ViewMyWallet(ec echo.Context) (err error) {
	var (
		ctx = ec.Request().Context()
		xid = tokenizer.FetchEchoTokenXID(ec)
	)

	wallet, cerr := c.UserWalletService.GetUserWalletByUserXID(ctx, xid)
	if cerr != nil {
		return ec.JSON(cerr.GetStatusCode(), jsend.GenerateResponseError(cerr))
	}

	return ec.JSON(http.StatusOK, jsend.GenerateResponseSuccess(map[string]interface{}{
		"wallet": map[string]interface{}{
			"id":         wallet.ID,
			"owned_by":   wallet.UserXID,
			"status":     wallet.IsEnabledString(),
			"enabled_at": wallet.EnabledAt,
			"balance":    wallet.CurrentBalance,
		},
	}))
}
