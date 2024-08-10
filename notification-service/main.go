package main

import (
	"log"
	"os"
	"strings"

	"github.com/0xivanov/crypto-notification-system/common/kafka"
	handler "github.com/0xivanov/crypto-notification-system/notification-service/consume_handler"
	"github.com/0xivanov/crypto-notification-system/notification-service/db"
	"github.com/0xivanov/crypto-notification-system/notification-service/notification"
)

func main() {
	// get brokers
	brokersString := os.Getenv("BROKERS")
	if brokersString == "" {
		brokersString = "localhost:9092"
	}
	brokers := strings.Split(brokersString, ",")

	// get topics
	topicsString := os.Getenv("TOPICS")
	if topicsString == "" {
		topicsString = "ticker,subscription,unsubscription"
	}
	topics := strings.Split(topicsString, ",")

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
	logger := log.New(os.Stdout, "notification-service", log.LstdFlags)

	// create notifiers
	mailNotifier := notification.NewMailNotifier(smtpHost, smtpUser, smtpPass, smtpUser, logger)
	slackNotifier := notification.NewSlackNotifier(logger)

	// create mongo client
	mongoClient := db.NewMongo(mongoURI, "notification-service", "users", logger)

	// create consumer for ticker updates
	tickerUpdateHandler := handler.NewTickerUpdateHandler(mongoClient, logger, []notification.Notifier{mailNotifier, slackNotifier})
	tickerConsumer := kafka.NewConsumer(brokers, topics[0], "notification-service", logger)
	go tickerConsumer.StartConsumer(tickerUpdateHandler)

	// create consumer for user subscriptions
	userSubscribeHandler := handler.NewUserSubscribeHandler(mongoClient, logger)
	subscriptionConsumer := kafka.NewConsumer(brokers, topics[1], "notification-service", logger)
	go subscriptionConsumer.StartConsumer(userSubscribeHandler)

	// create consumer for user unsubscriptions
	userUnsubscribeHandler := handler.NewUserUnsubscribeHandler(mongoClient, logger)
	unsubscriptionConsumer := kafka.NewConsumer(brokers, topics[2], "notification-service", logger)
	unsubscriptionConsumer.StartConsumer(userUnsubscribeHandler)
}
