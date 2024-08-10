package consume_handler

import (
	"encoding/json"
	"log"

	"github.com/0xivanov/crypto-notification-system/aggregator-service/kraken"
	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/IBM/sarama"
)

type UserSubscribeHandler struct {
	logger       *log.Logger
	krakenClient kraken.KrakenClientInterface
}

func NewUserSubscribeHandler(logger *log.Logger, krakenClient kraken.KrakenClientInterface) *UserSubscribeHandler {
	return &UserSubscribeHandler{
		logger:       logger,
		krakenClient: krakenClient,
	}
}

func (h *UserSubscribeHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *UserSubscribeHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *UserSubscribeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.handleMessage(message.Value)
		session.MarkMessage(message, "")
	}
	return nil
}

// handleMessage handles a message from the subscription topic
// and subscribes the user for the tickers
func (h *UserSubscribeHandler) handleMessage(message []byte) {
	var user model.User
	if err := json.Unmarshal(message, &user); err != nil {
		h.logger.Printf("[ERROR] Failed to unmarshal message: %v", err)
		return
	}

	for _, ticker := range user.Tickers {
		err := h.krakenClient.Subscribe(user.UserID, ticker.Symbol)
		if err != nil {
			h.logger.Printf("[INFO] Failed to subscribe to ticker %s: %v", ticker.Symbol, err)
			continue
		}
		h.logger.Printf("[INFO] UserID %s Subscribed to ticker %s", user.UserID, ticker.Symbol)
	}
}
