package consume_handler

import (
	"encoding/json"
	"log"

	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/0xivanov/crypto-notification-system/notification-service/db"
	"github.com/IBM/sarama"
)

type UserUnsubscribeHandler struct {
	db     *db.Mongo
	logger *log.Logger
}

func NewUserUnsubscribeHandler(db *db.Mongo, logger *log.Logger) *UserUnsubscribeHandler {
	return &UserUnsubscribeHandler{
		db:     db,
		logger: logger,
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

func (h *UserUnsubscribeHandler) handleMessage(message []byte) {
	var user model.User
	if err := json.Unmarshal(message, &user); err != nil {
		h.logger.Printf("[ERROR] Failed to unmarshal message: %v", err)
		return
	}

	if _, err := h.db.GetUserByID(user.UserID); err != nil {
		return
	}
	if err := h.db.RemoveUser(user.UserID); err != nil {
		return
	}
	log.Printf("[INFO] Unsubscribed user %s", user.UserID)
}
