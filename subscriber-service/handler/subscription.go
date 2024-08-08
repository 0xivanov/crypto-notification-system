package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	logger *log.Logger
}

func NewSubscriptionHandler(logger *log.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		logger: logger,
	}
}

func (s *SubscriptionHandler) Subscribe(ctx *gin.Context) {
	s.logger.Println("Subscribe handler")
}

func (s *SubscriptionHandler) Unsubscribe(ctx *gin.Context) {
	s.logger.Println("Unsubscribe handler")
}
