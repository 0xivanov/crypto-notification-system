package http_handler

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/0xivanov/crypto-notification-system/common/kafka"
	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/0xivanov/crypto-notification-system/subscriber-service/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	logger   *log.Logger
	producer kafka.ProducerInterface
}

func NewSubscriptionHandler(logger *log.Logger, producer kafka.ProducerInterface) *SubscriptionHandler {
	return &SubscriptionHandler{
		logger:   logger,
		producer: producer,
	}
}

func (s *SubscriptionHandler) Subscribe(ctx *gin.Context) {
	var user model.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.INVALID_REQUEST})
		return
	}
	if user.Tickers == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.INVALID_REQUEST})
		return
	}
	if user.NotificationOptions.Email == "" && user.NotificationOptions.PhoneNumber == "" && user.NotificationOptions.SlackWebhookURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.INVALID_REQUEST})
		return
	}
	if user.UserID == "" {
		user.UserID = uuid.New().String()
	}

	bytesUser, err := json.Marshal(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.INTERNAL_SERVER_ERROR})
		return
	}

	err = s.producer.SendMessage("subscription", string(bytesUser))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.INTERNAL_SERVER_ERROR})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "subscription request accepted, userID: " + user.UserID})
}

func (s *SubscriptionHandler) Unsubscribe(ctx *gin.Context) {
	var user model.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.INVALID_REQUEST})
		return
	}
	if user.UserID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.INVALID_REQUEST})
		return
	}
	bytesUser, err := json.Marshal(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.INTERNAL_SERVER_ERROR})
		return
	}

	err = s.producer.SendMessage("unsubscription", string(bytesUser))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.INTERNAL_SERVER_ERROR})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "unsubscription request accepted"})
}
