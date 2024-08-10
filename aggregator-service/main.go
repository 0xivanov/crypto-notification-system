package main

import (
	"log"
	"os"
	"strings"

	handler "github.com/0xivanov/crypto-notification-system/aggregator-service/consume_handler"
	"github.com/0xivanov/crypto-notification-system/aggregator-service/db"
	"github.com/0xivanov/crypto-notification-system/aggregator-service/kraken"
	"github.com/0xivanov/crypto-notification-system/common/kafka"
	"github.com/redis/go-redis/v9"
)

func main() {
	// get brokers
	brokersString := os.Getenv("BROKERS")
	if brokersString == "" {
		brokersString = "localhost:9092"
	}
	brokers := strings.Split(brokersString, ",")

	// get redis host
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost:6379"
	}

	// get redis password
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = ""
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

	// create redis cache
	redisDb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})
	redisCache := db.NewRedisCache(redisDb, logger)

	// create kraken client
	krakenClient := kraken.NewWebSocketClient(logger, producer, redisCache, wsUrl)

	// create subscribe consumer
	userSubscribeHandler := handler.NewUserSubscribeHandler(logger, krakenClient)
	subscribeConsumer := kafka.NewConsumer(brokers, "subscription", "aggregator-service", logger)
	go subscribeConsumer.StartConsumer(userSubscribeHandler)

	// create usubscribe consumer
	userUnsubscribeHandler := handler.NewUserUnsubscribeHandler(logger, krakenClient)
	unsubscribeConsumer := kafka.NewConsumer(brokers, "unsubscription", "aggregator-service", logger)
	go unsubscribeConsumer.StartConsumer(userUnsubscribeHandler)

	krakenClient.Listen()
}
