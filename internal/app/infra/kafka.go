package infra

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	kf "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type (
	// @envconfig (prefix:"BROKER_KAFKA")
	KafkaCfg struct {
		Address string `envconfig:"ADDRES" default:":9092"`
	}
)

func InitKafkaConfig(cfgApp *App, cfgKafka *KafkaCfg) *kf.ConfigMap {
	host, err := os.Hostname()
	if err != nil {
		host = "?"
	}

	return &kf.ConfigMap{
		"bootstrap.servers": cfgKafka.Address,
		"group.id":          fmt.Sprintf("%s-%s", cfgApp.Key, cfgApp.Env),
		"client.id":         host,
	}
}

// @ctor
func NewProducer(cfgApp *App, cfgKafka *KafkaCfg) *kf.Producer {
	producer, err := kf.NewProducer(InitKafkaConfig(cfgApp, cfgKafka))
	if err != nil {
		logrus.Fatalf("kafka producer: %s", err.Error())
	}
	logrus.Info("Connected to kafka producer")

	return producer
}

// @ctor
func NewConsumer(cfgApp *App, cfgKafka *KafkaCfg) *kf.Consumer {
	kafkaConsumer, err := kf.NewConsumer(InitKafkaConfig(cfgApp, cfgKafka))
	if err != nil {
		logrus.Fatalf("kafka consumer: %s", err.Error())
	}

	logrus.Info("Connected to kafka consumer")
	return kafkaConsumer
}
