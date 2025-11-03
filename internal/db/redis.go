package db

import (
	"context"
	"fmt"

	"github.com/dev-hyunsang/ticketly-backend/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	addr := config.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := config.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}
