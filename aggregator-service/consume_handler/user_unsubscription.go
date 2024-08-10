package consume_handler

import (
	"encoding/json"
	"log"

	"github.com/0xivanov/crypto-notification-system/aggregator-service/kraken"
	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/IBM/sarama"
)

type UserUnsubscribeHandler struct {
	logger       *log.Logger
	krakenClient kraken.KrakenClientInterface
}

func NewUserUnsubscribeHandler(logger *log.Logger, krakenClient kraken.KrakenClientInterface) *UserUnsubscribeHandler {
	return &UserUnsubscribeHandler{
		logger:       logger,
		krakenClient: krakenClient,
	}
}

func (h *UserUnsubscribeHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *UserUnsubscribeHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *UserUnsubscribeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.handleMessage(message.Value)
		session.MarkMessage(message, "")
	}
	return nil
}

// handleMessage handles a message from the unsubscription topic
// and unsubscribes the user for the tickers
func (h *UserUnsubscribeHandler) handleMessage(message []byte) {
	var user model.User
	if err := json.Unmarshal(message, &user); err != nil {
		h.logger.Printf("[ERROR] Failed to unmarshal message: %v", err)
		return
	}

	for _, ticker := range user.Tickers {
		err := h.krakenClient.Unsubscribe(user.UserID, ticker.Symbol)
		if err != nil {
			h.logger.Printf("[ERROR] Failed to unsubscribe user %s to ticker %s: %v", user.UserID, ticker.Symbol, err)
			return
		}
		h.logger.Printf("[INFO] Unsubscribed user %s to ticker %s", user.UserID, ticker.Symbol)
	}
}
