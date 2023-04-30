package kafka

import "github.com/freekup/mini-wallet/internal/app/entity"

type (
	MessageRepository interface {
		CreatedWalletTransaction(arg entity.KafkaCreatedWalletTransactionData) (err error)
	}
)
