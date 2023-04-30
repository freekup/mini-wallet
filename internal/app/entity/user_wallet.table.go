package entity

var (
	UserWalletTableName = "public.user_wallets"
	UserWalletTable     = struct {
		ID             string
		UserXID        string
		CurrentBalance string
		IsEnabled      string
		EnabledAt      string
		CreatedAt      string
		CreatedBy      string
		ModifiedAt     string
		ModifiedBy     string
		DeletedAt      string
		DeletedBy      string
	}{
		ID:             "id",
		UserXID:        "user_xid",
		CurrentBalance: "current_balance",
		IsEnabled:      "is_enabled",
		EnabledAt:      "enabled_at",
		CreatedAt:      "created_at",
		CreatedBy:      "created_by",
		ModifiedAt:     "modified_at",
		ModifiedBy:     "modified_by",
		DeletedAt:      "deleted_at",
		DeletedBy:      "deleted_by",
	}
)
