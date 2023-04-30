package entity

type (
	UserWallet struct {
		ID             string  `json:"id"`
		UserXID        string  `json:"user_xid"`
		CurrentBalance float64 `json:"balance"`
		IsEnabled      int     `json:"is_enabled"`
		EnabledAt      *string `json:"enabled_at"`
		DeletedBy      *string `json:"deleted_by"`
		DeletedAt      *string `json:"deleted_at"`
	}

	CreateUserWalletArg struct {
		UserXID string
	}
)

func (v UserWallet) StringIsEnabled() string {
	if v.IsEnabled == 1 {
		return "enabled"
	} else {
		return "disabled"
	}
}
