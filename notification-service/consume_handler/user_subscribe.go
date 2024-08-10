package consume_handler

import (
	"encoding/json"
	"log"

	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/0xivanov/crypto-notification-system/notification-service/db"
	"github.com/IBM/sarama"
)

type UserSubscribeHandler struct {
	db     *db.Mongo
	logger *log.Logger
}

func NewUserSubscribeHandler(db *db.Mongo, logger *log.Logger) *UserSubscribeHandler {
	return &UserSubscribeHandler{
		db:     db,
		logger: logger,
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

func (h *UserSubscribeHandler) handleMessage(message []byte) {
	var user model.User
	if err := json.Unmarshal(message, &user); err != nil {
		h.logger.Printf("[ERROR] Failed to unmarshal message: %v", err)
		return
	}
	if _, err := h.db.GetUserByID(user.UserID); err != nil {
		if err := h.db.AddUser(user); err != nil {
			return
		}
		log.Printf("[INFO] Subscribed user %s", user.UserID)
	} else {
		if err := h.db.UpdateUser(user.UserID, user); err != nil {
			return
		}
		log.Printf("[INFO] Updated user %s", user.UserID)
	}
}
