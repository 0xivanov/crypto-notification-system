package handler

import (
	"encoding/json"
	"log"

	"github.com/0xivanov/crypto-notification-system/notification-service/db"
	"github.com/0xivanov/crypto-notification-system/notification-service/model"
	"github.com/0xivanov/crypto-notification-system/notification-service/notification"
	"github.com/0xivanov/crypto-notification-system/notification-service/util"
	"github.com/IBM/sarama"
)

type TickerUpdateHandler struct {
	db        *db.Mongo
	notifiers []notification.Notifier
	logger    *log.Logger
}

func NewTickerUpdateHandler(db *db.Mongo, logger *log.Logger, notifiers []notification.Notifier) *TickerUpdateHandler {
	return &TickerUpdateHandler{
		db:        db,
		logger:    logger,
		notifiers: notifiers,
	}
}

func (h *TickerUpdateHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *TickerUpdateHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *TickerUpdateHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.handleMessage(message.Value)
		session.MarkMessage(message, "")
	}
	return nil
}

func (h *TickerUpdateHandler) handleMessage(message []byte) {
	var tickerUpdate model.Ticker
	if err := json.Unmarshal(message, &tickerUpdate); err != nil {
		h.logger.Printf("[ERROR] Failed to unmarshal message: %v", err)
		return
	}

	for _, tickerData := range tickerUpdate.Data {
		users, err := h.db.GetUsersForTicker(tickerData.Symbol)
		if err != nil {
			h.logger.Printf("[ERROR] Failed to get users for ticker %s: %v", tickerData.Symbol, err)
			continue
		}

		for _, user := range users {
			for _, notifier := range h.notifiers {
				err := notifier.SendNotification(util.FormatMessage(tickerData), user.NotificationOptions)
				if err != nil {
					h.logger.Printf("[ERROR] Failed to send notification to user %s: %v", user, err)
				}
			}
		}
	}
}
