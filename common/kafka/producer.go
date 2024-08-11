package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type ProducerInterface interface {
	SendMessage(topic, message string) error
}

type Producer struct {
	syncProducer sarama.SyncProducer
	logger       *log.Logger
}

func NewProducer(brokers []string, logger *log.Logger) *Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.Fatalf("[ERROR] Failed to start Kafka producer, connecting to broker %v", brokers)
	}

	return &Producer{
		syncProducer: producer,
		logger:       logger,
	}
}

func (p *Producer) SendMessage(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := p.syncProducer.SendMessage(msg)
	return err
}
