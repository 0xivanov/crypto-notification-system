package consume_handler

import (
	"encoding/json"
	"log"

	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/0xivanov/crypto-notification-system/notification-service/db"
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

// handleMessage handles a message from the ticker update topic
// and sends notifications to users if the ticker data reaches the threshold
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

		// loop through users
		for _, user := range users {
			// if does not reach threshold, skip
			if !isReachingThreshold(user, tickerData) {
				continue
			}

			// send notification through all notifiers
			for _, notifier := range h.notifiers {
				err := notifier.SendNotification(util.FormatMessage(tickerData), user.NotificationOptions)
				if err != nil {
					h.logger.Printf("[ERROR] Failed to send notification to user %s: %v", user.UserID, err)
				}
			}
		}
	}
}

func isReachingThreshold(user model.User, tickerData model.TickerData) bool {
	for _, ticker := range user.Tickers {
		if ticker.ChangeThreshold <= tickerData.ChangePct {
			return true
		}
	}
	return false
}
