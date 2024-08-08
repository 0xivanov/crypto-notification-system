package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/0xivanov/crypto-notification-system/subscriber-service/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "9090" // Default port if not provided
	}
	// create logger
	logger := log.New(os.Stdout, "notification-service ", log.LstdFlags)

	// create handlers
	subscriptionHandler := handler.NewSubscriptionHandler(logger)

	ginEngine := gin.Default()
	ginEngine.POST("/subscribe", subscriptionHandler.Subscribe)
	ginEngine.POST("/unsubscribe", subscriptionHandler.Unsubscribe)

	s := http.Server{
		Addr:         "0.0.0.0:" + port, // configure the bind address
		Handler:      ginEngine,         // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
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
	log.Println("[INFO] Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
