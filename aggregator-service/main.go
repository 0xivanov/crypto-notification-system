package main

import (
	"log"
	"os"
	"strings"

	"github.com/0xivanov/crypto-notification-system/aggregator-service/kafka"
	"github.com/0xivanov/crypto-notification-system/aggregator-service/kraken"
)

func main() {
	// get brokers
	brokersString := os.Getenv("BROKERS")
	if brokersString == "" {
		brokersString = "localhost:9092"
	}
	brokers := strings.Split(brokersString, ",")

	// get topic
	topic := os.Getenv("TOPIC")
	if topic == "" {
		topic = "ticker"
	}

	// get ws url
	wsUrl := os.Getenv("WS_URL")
	if wsUrl == "" {
		wsUrl = "wss://ws.kraken.com/v2"
	}

	// create logger
	logger := log.New(os.Stdout, "aggregator-service ", log.LstdFlags)

	// create producer
	producer := kafka.NewProducer(brokers, logger)

	// create consumer TODO

	krakenClient := kraken.NewWebSocketClient(logger, wsUrl, topic)
	krakenClient.Subscribe("BTC/USD")
	krakenClient.Subscribe("ETH/USD")
	krakenClient.Listen(producer)

}
