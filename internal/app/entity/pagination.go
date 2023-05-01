package entity

type (
	ViewPagination struct {
		Limit  int64 `json:"limit" query:"limit"`
		Offset int64 `json:"offset" query:"offset"`
		Total  int64 `json:"total"`
	}
)
