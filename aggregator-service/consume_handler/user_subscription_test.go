package consume_handler

import (
	"encoding/json"
	"errors"
	"log"
	"testing"

	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/stretchr/testify/assert"
)

func TestHandleMessageSub_Success(t *testing.T) {
	logger := log.Default()
	mockKrakenClient := new(MockKrakenClient)
	handler := NewUserSubscribeHandler(logger, mockKrakenClient)

	user := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
			{Symbol: "ETH/USD"},
		},
	}

	message, err := json.Marshal(user)
	assert.NoError(t, err)

	mockKrakenClient.On("Subscribe", "user123", "BTC/USD").Return(nil)
	mockKrakenClient.On("Subscribe", "user123", "ETH/USD").Return(nil)

	handler.handleMessage(message)

	mockKrakenClient.AssertExpectations(t)
}

func TestHandleMessage_SubscribeError(t *testing.T) {
	logger := log.Default()
	mockKrakenClient := new(MockKrakenClient)
	handler := NewUserSubscribeHandler(logger, mockKrakenClient)

	user := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
		},
	}

	message, err := json.Marshal(user)
	assert.NoError(t, err)

	mockKrakenClient.On("Subscribe", "user123", "BTC/USD").Return(errors.New("failed to subscribe"))

	handler.handleMessage(message)

	mockKrakenClient.AssertExpectations(t)
}
