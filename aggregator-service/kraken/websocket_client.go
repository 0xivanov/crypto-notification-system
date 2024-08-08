package kraken

import (
	"encoding/json"
	"log"

	"github.com/0xivanov/crypto-notification-system/aggregator-service/kafka"
	"github.com/0xivanov/crypto-notification-system/aggregator-service/model"
	"github.com/gorilla/websocket"
)

var subscriptionMessage = model.WebSocketMessage{
	Method: "subscribe",
	Params: model.Params{
		Channel: "ticker",
		Symbols: []string{},
	},
}

var unsubscriptionMessage = model.WebSocketMessage{
	Method: "unsubscribe",
	Params: model.Params{
		Channel: "ticker",
		Symbols: []string{},
	},
}

type WebSocketClient struct {
	socket     *websocket.Conn
	logger     *log.Logger
	kafkaTopic string
}

func NewWebSocketClient(logger *log.Logger, wsUrl, kafkaTopic string) *WebSocketClient {
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		logger.Fatalf("[ERROR] Failed to connect to web socket: %v", err)
	}

	return &WebSocketClient{
		socket:     conn,
		logger:     logger,
		kafkaTopic: kafkaTopic,
	}
}

func (c *WebSocketClient) Subscribe(ticker string) error {
	subscriptionMessage.Params.Symbols = append(subscriptionMessage.Params.Symbols, ticker)
	return c.socket.WriteJSON(subscriptionMessage)
}

func (c *WebSocketClient) Unsubscribe(ticker string) error {
	subscriptionMessage.Params.Symbols = append(subscriptionMessage.Params.Symbols, ticker)
	return c.socket.WriteJSON(unsubscriptionMessage)
}

func (c *WebSocketClient) Listen(producer *kafka.Producer) {
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
		for _, tickerData := range tickerUpdate.Data {
			c.logger.Printf("[INFO]: Received ticker update: %s - Last Price: %f", tickerData.Symbol, tickerData.Last)
			err := producer.SendMessage(c.kafkaTopic, string(message))
			if err != nil {
				log.Println(err)
			}
		}
	}
}
