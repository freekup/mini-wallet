package kafka

import (
	"encoding/json"

	"github.com/freekup/mini-wallet/internal/app/entity"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type (
	MessageRepositoryImpl struct {
		kafkaProducer *kafka.Producer
	}
)

// @ctor
func NewMessageRepository(kafkaProducer *kafka.Producer) MessageRepository {
	return &MessageRepositoryImpl{
		kafkaProducer: kafkaProducer,
	}
}

// CreatedWalletTransaction used to publish topic created wallet transaction
func (r *MessageRepositoryImpl) CreatedWalletTransaction(arg entity.KafkaCreatedWalletTransactionData) (err error) {
	return r.Publish(entity.KafkaTopicCreatedWalletTransaction, arg)
}

// Publish used to publish message
func (r *MessageRepositoryImpl) Publish(topic string, data interface{}) (err error) {
	byt, err := json.Marshal(data)
	if err != nil {
		return err
	}

	finEvn := make(chan kafka.Event)
	go func() {
		err = r.kafkaProducer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: byt,
		}, finEvn)
		if err != nil {
			finEvn <- nil
		}
	}()

	evn := <-finEvn
	if err != nil {
		return err
	}

	msg := evn.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		err = msg.TopicPartition.Error
		return err
	}

	return err
}
