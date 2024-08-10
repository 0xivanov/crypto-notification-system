package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	logger *log.Logger
}

func NewRedisCache(redisClient *redis.Client, logger *log.Logger) *RedisCache {
	return &RedisCache{
		client: redisClient,
		logger: logger,
	}
}

// AddUserForTicker adds a userID to the list of users subscribed to a ticker
func (r *RedisCache) AddUserForTicker(ticker string, userID string) error {
	ctx := context.Background()
	_, err := r.client.SAdd(ctx, ticker, userID).Result()
	if err != nil {
		r.logger.Printf("[ERROR] Failed to add user %s for ticker %s: %v", userID, ticker, err)
		return err
	}
	return nil
}

// RemoveUserForTicker removes a userID from the list of users subscribed to a ticker
func (r *RedisCache) RemoveUserForTicker(ticker string, userID string) error {
	ctx := context.Background()
	_, err := r.client.SRem(ctx, ticker, userID).Result()
	if err != nil {
		r.logger.Printf("[ERROR] Failed to remove user %s for ticker %s: %v", userID, ticker, err)
		return err
	}

	// if no users are left subscribed to the ticker, delete the ticker key
	membersCount, err := r.client.SCard(ctx, ticker).Result()
	if err != nil {
		r.logger.Printf("[ERROR] Failed to get remaining users count for ticker %s: %v", ticker, err)
		return err
	}
	if membersCount == 0 {
		err = r.client.Del(ctx, ticker).Err()
		if err != nil {
			r.logger.Printf("[ERROR] Failed to delete ticker %s after removing last user: %v", ticker, err)
			return err
		}
	}

	return nil
}

// GetUsersForTicker retrieves the list of userIDs subscribed to a ticker
func (r *RedisCache) GetUsersForTicker(ticker string) ([]string, int, error) {
	ctx := context.Background()
	userIDs, err := r.client.SMembers(ctx, ticker).Result()
	if err != nil {
		r.logger.Printf("[ERROR] Failed to get users for ticker %s: %v", ticker, err)
		return nil, 0, err
	}
	return userIDs, len(userIDs), nil
}

// GetAllTickers retrieves all tickers stored in the cache
func (r *RedisCache) GetAllTickers() ([]string, error) {
	ctx := context.Background()

	var (
		cursor uint64 // starting cursor position
		keys   []string
		err    error
	)

	// Iterate over the keys using SCAN with a match pattern for tickers.
	for {
		var scannedKeys []string
		scannedKeys, cursor, err = r.client.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			r.logger.Printf("[ERROR] Failed to retrieve tickers: %v", err)
			return nil, err
		}

		keys = append(keys, scannedKeys...)

		// If cursor is zero, the iteration is complete.
		if cursor == 0 {
			break
		}
	}

	return keys, nil
}
