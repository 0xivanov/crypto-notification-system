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

type MockKrakenClient struct {
	mock.Mock
}

func (m *MockKrakenClient) Subscribe(userID, ticker string) error {
	args := m.Called(userID, ticker)
	return args.Error(0)
}

func (m *MockKrakenClient) Unsubscribe(userID, ticker string) error {
	args := m.Called(userID, ticker)
	return args.Error(0)
}

func TestHandleMessageUnsub_Success(t *testing.T) {
	logger := log.Default()
	mockKrakenClient := new(MockKrakenClient)
	handler := NewUserUnsubscribeHandler(logger, mockKrakenClient)

	user := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
			{Symbol: "ETH/USD"},
		},
	}

	message, err := json.Marshal(user)
	assert.NoError(t, err)

	mockKrakenClient.On("Unsubscribe", "user123", "BTC/USD").Return(nil)
	mockKrakenClient.On("Unsubscribe", "user123", "ETH/USD").Return(nil)

	handler.handleMessage(message)

	mockKrakenClient.AssertExpectations(t)
}

func TestHandleMessage_UnsubscribeError(t *testing.T) {
	logger := log.Default()
	mockKrakenClient := new(MockKrakenClient)
	handler := NewUserUnsubscribeHandler(logger, mockKrakenClient)

	user := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
		},
	}

	message, err := json.Marshal(user)
	assert.NoError(t, err)

	mockKrakenClient.On("Unsubscribe", "user123", "BTC/USD").Return(errors.New("failed to unsubscribe"))

	handler.handleMessage(message)

	mockKrakenClient.AssertExpectations(t)
}
