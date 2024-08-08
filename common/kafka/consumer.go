package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	topic         string
	logger        *log.Logger
}

func NewConsumer(brokers []string, topic string, logger *log.Logger) *Consumer {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, "aggregator-service-group", config)
	if err != nil {
		logger.Fatalf("[ERROR] Failed to start Kafka consumer group: %v", err)
	}
	return &Consumer{
		consumerGroup: consumerGroup,
		topic:         topic,
		logger:        logger,
	}
}

func (c *Consumer) StartTickerUpdateConsumer(handler sarama.ConsumerGroupHandler) {
	for {
		err := c.consumerGroup.Consume(context.Background(), []string{c.topic}, handler)
		if err != nil {
			log.Fatalf("Error consuming Kafka messages: %v", err)
		}
	}
}
