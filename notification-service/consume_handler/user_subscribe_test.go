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

type MockMongoDB struct {
	mock.Mock
}

func (m *MockMongoDB) GetUserByID(userID string) (*model.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockMongoDB) AddUser(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockMongoDB) UpdateUser(userID string, user model.User) error {
	args := m.Called(userID, user)
	return args.Error(0)
}

func (m *MockMongoDB) RemoveUser(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockMongoDB) GetUsersForTicker(ticker string) ([]model.User, error) {
	args := m.Called(ticker)
	return args.Get(0).([]model.User), args.Error(2)
}

func TestHandleMessage_AddUser(t *testing.T) {
	logger := log.Default()
	mockDB := new(MockMongoDB)
	handler := NewUserSubscribeHandler(mockDB, logger)

	user := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
		},
	}

	message, err := json.Marshal(user)
	assert.NoError(t, err)

	mockDB.On("GetUserByID", "user123").Return(&model.User{}, errors.New("user not found"))
	mockDB.On("AddUser", user).Return(nil)

	handler.handleMessage(message)

	mockDB.AssertExpectations(t)
}

func TestHandleMessage_UpdateUser(t *testing.T) {
	logger := log.Default()
	mockDB := new(MockMongoDB)
	handler := NewUserSubscribeHandler(mockDB, logger)

	user := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
		},
	}

	message, err := json.Marshal(user)
	assert.NoError(t, err)

	mockDB.On("GetUserByID", "user123").Return(&user, nil)
	mockDB.On("UpdateUser", "user123", user).Return(nil)

	handler.handleMessage(message)

	mockDB.AssertExpectations(t)
}

func TestHandleMessage_AddUserError(t *testing.T) {
	logger := log.Default()
	mockDB := new(MockMongoDB)
	handler := NewUserSubscribeHandler(mockDB, logger)

	user := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
		},
	}

	message, err := json.Marshal(user)
	assert.NoError(t, err)

	mockDB.On("GetUserByID", "user123").Return(&model.User{}, errors.New("user not found"))
	mockDB.On("AddUser", user).Return(errors.New("failed to add user"))

	handler.handleMessage(message)

	mockDB.AssertExpectations(t)
}

func TestHandleMessage_UpdateUserError(t *testing.T) {
	logger := log.Default()
	mockDB := new(MockMongoDB)
	handler := NewUserSubscribeHandler(mockDB, logger)

	user := model.User{
		UserID: "user123",
		Tickers: []model.TickerSettings{
			{Symbol: "BTC/USD"},
		},
	}

	message, err := json.Marshal(user)
	assert.NoError(t, err)

	mockDB.On("GetUserByID", "user123").Return(&user, nil)
	mockDB.On("UpdateUser", "user123", user).Return(errors.New("failed to update user"))

	handler.handleMessage(message)

	mockDB.AssertExpectations(t)
}
