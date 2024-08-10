package consume_handler

import (
	"encoding/json"
	"errors"
	"log"
	"testing"

	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleMessage_UnsubscribeUser(t *testing.T) {
	logger := log.Default()
	mockDB := new(MockMongoDB)
	handler := NewUserUnsubscribeHandler(mockDB, logger)

	user := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
			{Symbol: "ETH/USD"},
		},
	}

	message, err := json.Marshal(user)
	assert.NoError(t, err)

	dbUser := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
			{Symbol: "ETH/USD"},
			{Symbol: "LTC/USD"},
		},
	}

	// Set up mock expectations
	mockDB.On("GetUserByID", "user123").Return(&dbUser, nil)
	mockDB.On("UpdateUser", "user123", model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "LTC/USD"},
		},
	}).Return(nil)

	handler.handleMessage(message)

	mockDB.AssertExpectations(t)
}

func TestHandleMessage_UserNotFound(t *testing.T) {
	logger := log.Default()
	mockDB := new(MockMongoDB)
	handler := NewUserUnsubscribeHandler(mockDB, logger)

	user := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
		},
	}

	message, err := json.Marshal(user)
	assert.NoError(t, err)

	mockDB.On("GetUserByID", "user123").Return(&model.User{}, errors.New("user not found"))

	handler.handleMessage(message)

	mockDB.AssertExpectations(t)
}

func TestHandleMessageUnsub_UpdateUserError(t *testing.T) {
	logger := log.Default()
	mockDB := new(MockMongoDB)
	handler := NewUserUnsubscribeHandler(mockDB, logger)

	user := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
		},
	}

	message, err := json.Marshal(user)
	assert.NoError(t, err)

	dbUser := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
			{Symbol: "ETH/USD"},
		},
	}

	mockDB.On("GetUserByID", "user123").Return(&dbUser, nil)
	mockDB.On("UpdateUser", "user123", mock.Anything).Return(errors.New("failed to update user"))

	handler.handleMessage(message)

	mockDB.AssertExpectations(t)
}
