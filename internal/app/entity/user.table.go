package entity

var (
	UserTableName = "public.users"
	UserTable     = struct {
		ID         string
		Name       string
		XID        string
		IsActive   string
		CreatedAt  string
		CreatedBy  string
		ModifiedAt string
		ModifiedBy string
		DeletedAt  string
		DeletedBy  string
	}{
		ID:         "id",
		Name:       "name",
		XID:        "xid",
		IsActive:   "is_active",
		CreatedAt:  "created_at",
		CreatedBy:  "created_by",
		ModifiedAt: "modified_at",
		ModifiedBy: "modified_by",
		DeletedAt:  "deleted_at",
		DeletedBy:  "deleted_by",
	}
)
