package redis 

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
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	RPush(ctx context.Context, key string, value interface{}) *redis.IntCmd
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
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

/* deletes a redis entry */
func (r *redisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, keys...)
}

/* for pushing multiple elements in Redis */
func (r *redisClient) RPush(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	return r.client.RPush(ctx, key, value)
}

/* retrieve a subset of the list stored at a specified key */
func (r *redisClient) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return r.client.LRange(ctx, key, start, stop)
}

/* hash set for redis */
func (r *redisClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return r.client.HSet(ctx, key, values...)
}
