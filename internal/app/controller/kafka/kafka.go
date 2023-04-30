package kafka

import (
	"context"
	"fmt"
	"os"

	"github.com/freekup/mini-wallet/internal/app/entity"
	uws "github.com/freekup/mini-wallet/internal/app/service/user_wallet"
	"github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/pkg/logruskit"
	"go.uber.org/dig"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type handler func(context.Context, *kafka.Message) error

type KafkaCtrl struct {
	dig.In
	KafkaConsumer     *kafka.Consumer
	Handlers          map[string]handler `optional:"true"`
	UserWalletService uws.UserWalletService
}

func (ox *KafkaCtrl) KafkaRoute() {
	ox.register(entity.KafkaTopicCreatedWalletTransaction, ox.CreatedWalletTransaction)

	_ = ox.register
	if err := ox.KafkaConsumer.SubscribeTopics(ox.topics(), nil); err != nil {
		logrus.Errorf("while SubscribeTopics detail: %s", err)
		return
	}

	go func() {
		for {
			msg, err := ox.KafkaConsumer.ReadMessage(-1)
			if err != nil {
				logrus.Errorf("kafka consumer err:%v", err)
				continue
			}

			ox.handleMessage(msg)
		}
	}()
}

func (ox *KafkaCtrl) register(topic string, handlr handler) (err error) {
	if ox.Handlers == nil {
		ox.Handlers = make(map[string]handler)
	}

	if _, ok := ox.Handlers[topic]; ok {
		err = fmt.Errorf("topic %s already exists", topic)
		return
	}

	ox.Handlers[topic] = handlr
	return
}

func (ox KafkaCtrl) topics() (topics []string) {
	for k := range ox.Handlers {
		topics = append(topics, k)
	}
	return topics
}

func (ox KafkaCtrl) handleMessage(msg *kafka.Message) {
	defer func() {
		if errs := recover(); errs != nil {
			logrus.Errorf("kafka consumer err:%v", errs)
		}
	}()

	topic := ""
	if msg.TopicPartition.Topic != nil {
		topic = *msg.TopicPartition.Topic
	}

	if handlr, ok := ox.Handlers[topic]; ok {
		ctx := context.Background()

		host, err := os.Hostname()
		if err != nil {
			host = "?"
		}

		logruskit.PutField(&ctx, "hostname", host)

		err = handlr(ctx, msg)
		if err != nil {
			logrus.Error(err)
		}

		return
	}
}
