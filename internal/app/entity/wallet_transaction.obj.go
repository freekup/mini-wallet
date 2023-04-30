package entity

type (
	WalletTransaction struct {
		ID          string  `json:"id"`
		WalletID    string  `json:"wallet_id"`
		Status      string  `json:"-"`
		ReferenceID string  `json:"reference_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		CreatedBy   string  `json:"created_by"`
		CreatedAt   string  `json:"created_at"`
	}
)

type (
	AddBalanceWalletArg struct {
		ReferenceID string
		Amount      float64
		Requestor   string
	}

	WithdrawBalanceArg struct {
		ReferenceID string
		Amount      float64
		Requestor   string
	}
)
