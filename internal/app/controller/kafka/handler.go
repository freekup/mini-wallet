package kafka

import (
	"context"
	"encoding/json"

	"github.com/freekup/mini-wallet/internal/app/entity"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// CreatedWalletTransaction used to handle incoming message from created wallet transaction
func (ox *KafkaCtrl) CreatedWalletTransaction(ctx context.Context, msg *kafka.Message) (err error) {
	param := entity.KafkaCreatedWalletTransactionData{}
	err = json.Unmarshal(msg.Value, &param)
	if err != nil {
		return err
	}

	return ox.UserWalletService.RefreshUserWalletCache(ctx, param.UserXID)
}
