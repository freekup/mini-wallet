package entity

var (
	WalletTransactionTableName = "public.wallet_transactions"
	WalletTransactionTable     = struct {
		ID          string
		WalletID    string
		Status      string
		ReferenceID string
		Amount      string
		Description string
		IsActive    string
		CreatedAt   string
		CreatedBy   string
		ModifiedAt  string
		ModifiedBy  string
		DeletedAt   string
		DeletedBy   string
	}{
		ID:          "id",
		WalletID:    "wallet_id",
		Status:      "status",
		ReferenceID: "reference_id",
		Amount:      "amount",
		Description: "description",
		IsActive:    "is_active",
		CreatedAt:   "created_at",
		CreatedBy:   "created_by",
		ModifiedAt:  "modified_at",
		ModifiedBy:  "modified_by",
		DeletedAt:   "deleted_at",
		DeletedBy:   "deleted_by",
	}
)
