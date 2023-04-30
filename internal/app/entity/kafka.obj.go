package entity

type (
	KafkaCreatedWalletTransactionData struct {
		UserXID string  `json:"user_xid"`
		Amount  float64 `json:"amount"`
	}
)
