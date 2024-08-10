package http_handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/google/uuid"
)

// MockKafkaProducer is a mock implementation of kafka.Producer
type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) SendMessage(topic string, message string) error {
	args := m.Called(topic, message)
	return args.Error(0)
}

// TestValidSubscription tests the Subscribe method with a valid subscription request
func TestValidSubscription(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	logger := log.Default()
	handler := NewSubscriptionHandler(logger, mockProducer)

	router := gin.Default()
	router.POST("/subscribe", handler.Subscribe)

	user := model.User{
		Tickers: []model.TickerSettings{
			{
				Symbol:          "BTC-USD",
				ChangeThreshold: 5.0,
			},
			{
				Symbol:          "ETH-USD",
				ChangeThreshold: 2.5,
			},
		},
		NotificationOptions: model.NotificationOptions{
			SlackWebhookURL: "https://hooks.slack.com/services/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",
			Email:           "test@example.com",
			PhoneNumber:     "+1234567890",
		},
	}

	mockProducer.On("SendMessage", "subscription", mock.Anything).Return(nil).Once()

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "subscription request accepted")
	mockProducer.AssertExpectations(t)
}

// TestInvalidSubscriptionWithoutTickers tests the Subscribe method with an invalid subscription request (no tickers provided)
func TestInvalidSubscriptionWithoutTickers(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	logger := log.Default()
	handler := NewSubscriptionHandler(logger, mockProducer)

	router := gin.Default()
	router.POST("/subscribe", handler.Subscribe)

	user := model.User{
		NotificationOptions: model.NotificationOptions{
			SlackWebhookURL: "https://hooks.slack.com/services/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",
			Email:           "test@example.com",
			PhoneNumber:     "+1234567890",
		},
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Invalid request")
	mockProducer.AssertNotCalled(t, "SendMessage", "subscription", mock.Anything)
}

// TestInvalidSubscriptionWithoutNotificationOptions tests the Subscribe method with an invalid subscription request (no notification options provided)
func TestInvalidSubscriptionWithoutNotificationOptions(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	logger := log.Default()
	handler := NewSubscriptionHandler(logger, mockProducer)

	router := gin.Default()
	router.POST("/subscribe", handler.Subscribe)

	user := model.User{
		Tickers: []model.TickerSettings{
			{
				Symbol:          "BTC-USD",
				ChangeThreshold: 5.0,
			},
		},
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Invalid request")
	mockProducer.AssertNotCalled(t, "SendMessage", "subscription", mock.Anything)
}

// TestValidUnsubscription tests the Unsubscribe method with a valid unsubscription request
func TestValidUnsubscription(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	logger := log.Default()
	handler := NewSubscriptionHandler(logger, mockProducer)

	router := gin.Default()
	router.POST("/unsubscribe", handler.Unsubscribe)

	user := model.User{
		UserID: uuid.New().String(),
	}

	mockProducer.On("SendMessage", "unsubscription", mock.Anything).Return(nil).Once()

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/unsubscribe", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "unsubscription request accepted")
	mockProducer.AssertExpectations(t)
}

// TestInvalidUnsubscriptionWithoutUserID tests the Unsubscribe method with an invalid unsubscription request (UserID not provided)
func TestInvalidUnsubscriptionWithoutUserID(t *testing.T) {
	mockProducer := new(MockKafkaProducer)
	logger := log.Default()
	handler := NewSubscriptionHandler(logger, mockProducer)

	router := gin.Default()
	router.POST("/unsubscribe", handler.Unsubscribe)

	user := model.User{
		UserID: "",
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/unsubscribe", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Invalid request")
	mockProducer.AssertNotCalled(t, "SendMessage", "unsubscription", mock.Anything)
}
