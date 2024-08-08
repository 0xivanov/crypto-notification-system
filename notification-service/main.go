package main

import (
	"log"
	"os"
	"strings"

	"github.com/0xivanov/crypto-notification-system/notification-service/db"
	"github.com/0xivanov/crypto-notification-system/notification-service/kafka"
	"github.com/0xivanov/crypto-notification-system/notification-service/kafka/handler"
	"github.com/0xivanov/crypto-notification-system/notification-service/notification"
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

	// get mongo uri
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// get smtp options
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	// create logger
	logger := log.New(os.Stdout, "notification-service ", log.LstdFlags)

	consumer := kafka.NewConsumer(brokers, topic, logger)
	mailNotifier := notification.NewMailNotifier(smtpHost, smtpUser, smtpPass, smtpUser, logger)
	slackNotifier := notification.NewSlackNotifier(logger)

	mongoClient := db.NewMongo(mongoURI, "notification-service", "users", logger)

	tickerUpdateHandler := handler.NewTickerUpdateHandler(mongoClient, logger, []notification.Notifier{mailNotifier, slackNotifier})
	consumer.StartTickerUpdateConsumer(tickerUpdateHandler)
}
