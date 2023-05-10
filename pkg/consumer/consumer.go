package consumer

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Consumer struct {
	logger *logrus.Logger
	sarama.PartitionConsumer
}

func NewConsumer(config viper.Viper, logger *logrus.Logger) (*Consumer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Return.Errors = true

	kafkaAddr := fmt.Sprintf("%s:%s", config.GetString("kafka.host"), config.GetString("kafka.port"))

	consumer, err := sarama.NewConsumer([]string{kafkaAddr}, kafkaConfig)
	if err != nil {
		return nil, err
	}

	topic := config.GetString("kafka.topic")
	partitionStr := config.GetString("kafka.partiton")
	partition, err := strconv.ParseInt(partitionStr, 10, 16)

	consumerPartition, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		logger: logger,
		PartitionConsumer: consumerPartition,
	}, nil
}

func (cons *Consumer) RunListener(userIdChan chan int, errChan chan error, signalChan chan os.Signal) {
	for {
		select {
		case err := <-cons.Errors():
			errChan <- err
			return
		case msg := <-cons.Messages():
			userID, err := strconv.ParseInt(string(msg.Value), 10, 32)
			if err != nil {
				errChan <- err
				return
			}

			if msg.Key[0] == 'R' {
				cons.logger.Infof("KAFKA MSG: KEY - %s, VALUE - %s", msg, string(msg.Key), string(msg.Value))
				userIdChan <- int(userID)
			}

		case <-signalChan:
			cons.logger.Infoln("Shutting down Kafka consumer...")
			cons.Close()
			close(userIdChan)
			return
		}
	}
}