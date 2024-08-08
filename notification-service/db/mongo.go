package db

import (
	"context"
	"log"
	"time"

	"github.com/0xivanov/crypto-notification-system/notification-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client     *mongo.Client
	dbName     string
	collection string
	logger     *log.Logger
}

func NewMongo(uri, dbName, collection string, logger *log.Logger) *Mongo {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Fatalf("[ERROR] Failed to connect to MongoDB: %v", err)
	}
	return &Mongo{
		client:     client,
		dbName:     dbName,
		collection: collection,
		logger:     logger,
	}
}

func (mongo *Mongo) GetUsersForTicker(ticker string) ([]model.User, error) {
	collection := mongo.client.Database(mongo.dbName).Collection(mongo.collection)
	filter := bson.M{"tickers": ticker}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var users []model.User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	mongo.logger.Printf("[INFO]: Found %d users for ticker %s", len(users), ticker)
	return users, nil
}
