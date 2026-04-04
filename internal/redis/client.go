package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/viral32111/glow-cad-exporter/internal/config"
)

type Client struct {
	Client    *redis.Client
	KeyPrefix string
}

func NewClient(context context.Context, config config.Config) (*Client, error) {
	redisOptions, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse Redis URL '%s': %w", config.RedisURL, err)
	}
	slog.Debug("Parsed Redis URL", "url", config.RedisURL)

	redisClient := redis.NewClient(redisOptions)
	slog.Info("Opened Redis connection", "address", redisOptions.Addr, "user", redisOptions.Username, "db", redisOptions.DB)

	if err := redisClient.Ping(context).Err(); err != nil {
		return nil, fmt.Errorf("Unable to ping Redis at '%s': %w", redisOptions.Addr, err)
	}
	slog.Debug("Pinged Redis", "address", redisOptions.Addr, "user", redisOptions.Username, "db", redisOptions.DB)

	return &Client{
		Client:    redisClient,
		KeyPrefix: config.RedisKeyPrefix,
	}, nil
}

func (client *Client) Close() error {
	slog.Info("Closing Redis connection", "address", client.Client.Options().Addr)
	return client.Client.Close()
}

func (client *Client) Set(context context.Context, key string, value string, ttl time.Duration) error {
	redisKey := client.KeyPrefix + key

	slog.Debug("Setting Redis key", "key", redisKey, "value", value)

	return client.Client.Set(context, redisKey, value, ttl).Err()
}

func (client *Client) Get(context context.Context, key string) (string, error) {
	redisKey := client.KeyPrefix + key

	slog.Debug("Getting Redis key", "key", redisKey)

	value, err := client.Client.Get(context, redisKey).Result()
	if err == redis.Nil {
		slog.Warn("Redis key not found", "key", redisKey)
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("Error getting Redis key '%s': %w", redisKey, err)
	}

	return value, nil
}
