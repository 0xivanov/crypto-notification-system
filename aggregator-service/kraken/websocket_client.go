package kraken

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/0xivanov/crypto-notification-system/aggregator-service/db"
	"github.com/0xivanov/crypto-notification-system/common/kafka"
	"github.com/0xivanov/crypto-notification-system/common/model"
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	socket   *websocket.Conn
	logger   *log.Logger
	producer kafka.ProducerInterface
	redis    *db.RedisCache
}

func NewWebSocketClient(logger *log.Logger, producer kafka.ProducerInterface, redis *db.RedisCache, wsUrl string) *WebSocketClient {
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		logger.Fatalf("[ERROR] Failed to connect to web socket: %v", err)
	}

	// get all tickers from Redis and subscribe to them
	tickers, err := redis.GetAllTickers()
	if err != nil {
		logger.Fatalf("[ERROR] Failed to get tickers from Redis: %v", err)
	}

	// subscribe to all tickers that exist in redis
	var subscriptionMessage = model.WebSocketMessage{
		Method: "subscribe",
		Params: model.Params{
			Channel: "ticker",
			Symbols: []string{},
		},
	}
	subscriptionMessage.Params.Symbols = append(subscriptionMessage.Params.Symbols, tickers...)
	err = conn.WriteJSON(subscriptionMessage)
	if err != nil {
		logger.Fatalf("[ERROR] Failed to send initial subscription message: %v", err)
	}

	return &WebSocketClient{
		socket:   conn,
		logger:   logger,
		producer: producer,
		redis:    redis,
	}
}

// Subscribe adds a user to the list of subscribers for a ticker
// and sends a subscription message if it's the first subscriber
func (c *WebSocketClient) Subscribe(userID, ticker string) error {
	// if the user is already subscribed - return
	userIDs, count, _ := c.redis.GetUsersForTicker(ticker)
	for _, id := range userIDs {
		if id == userID {
			return errors.New("user is already subscribed")
		}
	}

	// remove the user for this ticker
	err := c.redis.AddUserForTicker(ticker, userID)
	if err != nil {
		return err
	}

	// send WebSocket subscription message if it's the first subscriber
	var subscriptionMessage = model.WebSocketMessage{
		Method: "subscribe",
		Params: model.Params{
			Channel: "ticker",
			Symbols: []string{""},
		},
	}
	if count == 0 {
		subscriptionMessage.Params.Symbols[0] = ticker
		return c.socket.WriteJSON(subscriptionMessage)
	}

	return nil
}

// Unsubscribe removes a user from the list of subscribers for a ticker
// and sends an unsubscription message if there are no more subscribers
func (c *WebSocketClient) Unsubscribe(userID, ticker string) error {
	// if the user is not subscribed - return
	userIDs, count, _ := c.redis.GetUsersForTicker(ticker)
	fmt.Println(userIDs)
	isSubbed := false
	for _, id := range userIDs {
		if id == userID {
			isSubbed = true
		}
	}
	if !isSubbed {
		return errors.New("user is not subscribed")
	}

	// remove the user for this ticker
	err := c.redis.RemoveUserForTicker(ticker, userID)
	if err != nil {
		return err
	}

	// send WebSocket unsubscription message if there are no more subscribers
	var unsubscriptionMessage = model.WebSocketMessage{
		Method: "unsubscribe",
		Params: model.Params{
			Channel: "ticker",
			Symbols: []string{""},
		},
	}
	if count == 1 {
		unsubscriptionMessage.Params.Symbols[0] = ticker
		return c.socket.WriteJSON(unsubscriptionMessage)
	}

	return nil
}

// Listen listens for incoming messages from the WebSocket connection
func (c *WebSocketClient) Listen() {
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			c.logger.Printf("[ERROR] Failed to read message: %v", err)
			continue
		}
		var tickerUpdate model.Ticker
		if err := json.Unmarshal(message, &tickerUpdate); err != nil {
			c.logger.Printf("[ERROR] Failed to unmarshal ticker message: %v", err)
			continue
		}
		if tickerUpdate.Channel == "status" || tickerUpdate.Channel == "heartbeat" {
			continue
		}
		for _, tickerData := range tickerUpdate.Data {
			c.logger.Printf("[INFO]: Received ticker update: %s - Last Price: %f", tickerData.Symbol, tickerData.Last)
			err := c.producer.SendMessage("ticker", string(message))
			if err != nil {
				c.logger.Printf("[ERROR] Failed to send message to Kafka: %v", err)
			}
			c.logger.Printf("[INFO] Ticker symbol {%s} update sent to topic {%s}:", tickerData.Symbol, "ticker")
		}
	}
}
