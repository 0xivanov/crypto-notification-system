package consume_handler

import (
	"encoding/json"
	"log"

	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/0xivanov/crypto-notification-system/notification-service/db"
	"github.com/IBM/sarama"
)

type UserUnsubscribeHandler struct {
	db     db.MongoInterface
	logger *log.Logger
}

func NewUserUnsubscribeHandler(db db.MongoInterface, logger *log.Logger) *UserUnsubscribeHandler {
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

// handleMessage handles a message from the unsubscription topic
// and deletes the user from the db if
func (h *UserUnsubscribeHandler) handleMessage(message []byte) {
	var user model.User
	if err := json.Unmarshal(message, &user); err != nil {
		h.logger.Printf("[ERROR] Failed to unmarshal message: %v", err)
		return
	}

	// get the user
	dbUser, err := h.db.GetUserByID(user.UserID)
	if err != nil {
		return
	}

	// remove the tickers from the user's preferences
	for _, ticker := range user.Tickers {
		dbUser.Tickers = removeTicker(ticker.Symbol, dbUser.Tickers)
	}
	user.Tickers = dbUser.Tickers

	if err := h.db.UpdateUser(user.UserID, user); err != nil {
		return
	}
	log.Printf("[INFO] Unsubscribed user's %s tickers %v", user.UserID, user.Tickers)
}

func removeTicker(ticker string, tickers []model.TickerSettings) []model.TickerSettings {

	for i, t := range tickers {
		if t.Symbol == ticker {
			tickers = append(tickers[:i], tickers[i+1:]...)
		}
	}
	return tickers
}
