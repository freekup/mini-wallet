package rest

import (
	"github.com/freekup/mini-wallet/internal/app/entity"
	"github.com/freekup/mini-wallet/pkg/tokenizer"
	"github.com/labstack/echo/v4"
)

// GetWalletTransactions used to handle get list of transactions
func (c *UserWalletController) GetWalletTransactions(ec echo.Context) (err error) {
	var (
		ctx        = ec.Request().Context()
		xid        = tokenizer.FetchEchoTokenXID(ec)
		pagination = entity.ViewPagination{}
	)

	if err = ec.Bind(&pagination); err != nil {
		return
	}

	results, pag, err := c.WalletTransService.GetWalletTransactions(ctx, pagination, xid)
	if err != nil {
		return
	}

	return ec.JSON(200, map[string]interface{}{
		"data":       results,
		"pagination": pag,
	})
}

// WalletDeposit used to handle deposit
func (c *UserWalletController) WalletDeposit(ec echo.Context) (err error) {
	var (
		ctx   = ec.Request().Context()
		xid   = tokenizer.FetchEchoTokenXID(ec)
		param = struct {
			Amount      float64 `form:"amount"`
			ReferenceID string  `form:"reference_id"`
		}{}
	)

	if err = ec.Bind(&param); err != nil {
		return
	}

	walletTrans, err := c.WalletTransService.AddBalanceWallet(ctx, entity.AddBalanceWalletArg{
		Amount:      param.Amount,
		ReferenceID: param.ReferenceID,
		Requestor:   xid,
	})
	if err != nil {
		return
	}

	return ec.JSON(200, map[string]interface{}{
		"id":           walletTrans.ID,
		"deposit_by":   walletTrans.CreatedBy,
		"status":       walletTrans.Status,
		"deposit_at":   walletTrans.CreatedAt,
		"amount":       walletTrans.Amount,
		"reference_id": walletTrans.ReferenceID,
	})
}

// WalletWithdraw used to handle withdrawal
func (c *UserWalletController) WalletWithdraw(ec echo.Context) (err error) {
	var (
		ctx   = ec.Request().Context()
		xid   = tokenizer.FetchEchoTokenXID(ec)
		param = struct {
			Amount      float64 `form:"amount"`
			ReferenceID string  `form:"reference_id"`
		}{}
	)

	if err = ec.Bind(&param); err != nil {
		return
	}

	walletTrans, err := c.WalletTransService.WithdrawBalance(ctx, entity.WithdrawBalanceArg{
		Amount:      -param.Amount,
		ReferenceID: param.ReferenceID,
		Requestor:   xid,
	})
	if err != nil {
		return
	}

	return ec.JSON(200, map[string]interface{}{
		"id":           walletTrans.ID,
		"deposit_by":   walletTrans.CreatedBy,
		"status":       walletTrans.Status,
		"deposit_at":   walletTrans.CreatedAt,
		"amount":       walletTrans.Amount,
		"reference_id": walletTrans.ReferenceID,
	})
}
