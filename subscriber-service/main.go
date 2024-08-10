package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/0xivanov/crypto-notification-system/common/kafka"
	handler "github.com/0xivanov/crypto-notification-system/subscriber-service/http_handler"
	"github.com/gin-gonic/gin"
)

func main() {
	// get brokers
	brokersString := os.Getenv("BROKERS")
	if brokersString == "" {
		brokersString = "localhost:9092"
	}
	brokers := strings.Split(brokersString, ",")

	// get port
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	// create logger
	logger := log.New(os.Stdout, "notification-service", log.LstdFlags)

	// create producer
	producer := kafka.NewProducer(brokers, logger)

	// create http handlers
	subscriptionHandler := handler.NewSubscriptionHandler(logger, producer)

	ginEngine := gin.Default()
	ginEngine.POST("/subscribe", subscriptionHandler.Subscribe)
	ginEngine.POST("/unsubscribe", subscriptionHandler.Unsubscribe)

	s := http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      ginEngine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		logger.Printf("[INFO] Starting server on port %s", port)

		err := s.ListenAndServe()
		if err != nil {
			logger.Printf("[INFO] Closing server %d", err)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	logger.Println("[INFO] Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
