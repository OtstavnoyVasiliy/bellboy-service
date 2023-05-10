package producer

import (
	"encoding/json"
	"fmt"
	"tg-bot/pkg/types"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

type Producer struct {
	sarama.SyncProducer
	Topic string
}

func NewProducer(config viper.Viper) (*Producer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Errors = true
	kafkaConfig.Producer.Return.Successes = true

	kafkaAddr := fmt.Sprintf("%s:%s", config.GetString("kafka.host"), config.GetString("kafka.port"))
	producer, err := sarama.NewSyncProducer([]string{kafkaAddr}, kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &Producer{
		Topic: config.GetString("kafka.topic"),
		SyncProducer: producer,
	}, nil
}

func (p *Producer) SendKafkaMessage(msgVaue types.KickMessage) error {
	strBytes, err := json.Marshal(msgVaue)
	if err != nil {
		return err
	}

	data := sarama.StringEncoder(strBytes)

	msg := &sarama.ProducerMessage {
		Topic: p.Topic,
		Value: data,
	}

	if _, _, err := p.SendMessage(msg); err != nil {
		return err
	}

	return nil
}
