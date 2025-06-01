package session

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

/* defines the methods to expose (for dependency injection) */
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	HMSet(ctx context.Context, key string, values map[string]interface{}) error
}

/* redisClient implementation */
type redisClient struct {
	client *redis.Client
}

/* for creating a new redis client */
func NewRedisClient(address, password string, db int) (RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}

	return &redisClient{client: rdb}, nil
}

/* Set sets a key-value pair in Redis */
func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

/* retrieves a value by key from Redis */
func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

/* avoid unsafe casting */
func (r *redisClient) HMSet(ctx context.Context, key string, values map[string]interface{}) error {
	return r.client.HSet(ctx, key, values).Err()
}
