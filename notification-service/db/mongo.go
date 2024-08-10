package db

import (
	"context"
	"log"
	"time"

	"github.com/0xivanov/crypto-notification-system/common/model"
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

func (m *Mongo) GetUsersForTicker(symbol string) ([]model.User, error) {
	collection := m.client.Database(m.dbName).Collection(m.collection)

	filter := bson.M{
		"tickers": bson.M{
			"$elemMatch": bson.M{
				"symbol": symbol,
			},
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var users []model.User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}

	m.logger.Printf("[INFO]: Found %d users for ticker %s", len(users), symbol)
	return users, nil
}

func (m *Mongo) AddUser(user model.User) error {
	collection := m.client.Database(m.dbName).Collection(m.collection)

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		m.logger.Printf("[ERROR] Failed to insert user: %v", err)
		return err
	}

	m.logger.Printf("[INFO] User %s added successfully", user.UserID)
	return nil
}

// RemoveUser removes a user by their userID
func (m *Mongo) RemoveUser(userID string) error {
	collection := m.client.Database(m.dbName).Collection(m.collection)

	filter := bson.M{"userID": userID}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		m.logger.Printf("[ERROR] Failed to remove user %s: %v", userID, err)
		return err
	}

	if result.DeletedCount == 0 {
		m.logger.Printf("[INFO] No user found with userID %s to remove", userID)
	} else {
		m.logger.Printf("[INFO] User %s removed successfully", userID)
	}

	return nil
}

// UpdateUser updates the information of an existing user by their userID
func (m *Mongo) UpdateUser(userID string, updatedUser model.User) error {
	collection := m.client.Database(m.dbName).Collection(m.collection)

	filter := bson.M{"userID": userID}
	update := bson.M{"$set": updatedUser}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		m.logger.Printf("[ERROR] Failed to update user %s: %v", userID, err)
		return err
	}

	if result.MatchedCount == 0 {
		m.logger.Printf("[INFO] No user found with userID %s to update", userID)
	} else {
		m.logger.Printf("[INFO] User %s updated successfully", userID)
	}

	return nil
}

// GetUserByID retrieves a user by their userID
func (m *Mongo) GetUserByID(userID string) (*model.User, error) {
	collection := m.client.Database(m.dbName).Collection(m.collection)

	filter := bson.M{"userID": userID}

	var user model.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Printf("[INFO] No user found with userID %s", userID)
			return nil, err
		}
		m.logger.Printf("[ERROR] Failed to get user %s: %v", userID, err)
		return nil, err
	}

	m.logger.Printf("[INFO] User %s retrieved successfully", userID)
	return &user, nil
}
