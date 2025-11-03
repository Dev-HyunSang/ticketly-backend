package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) *TokenRepository {
	return &TokenRepository{
		client: client,
	}
}

// SaveRefreshToken saves refresh token with expiration
func (r *TokenRepository) SaveRefreshToken(userID uuid.UUID, token string, expiration time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%s", userID.String())

	err := r.client.Set(ctx, key, token, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

// GetRefreshToken retrieves refresh token by user ID
func (r *TokenRepository) GetRefreshToken(userID uuid.UUID) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%s", userID.String())

	token, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("refresh token not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}

	return token, nil
}

// DeleteRefreshToken removes refresh token
func (r *TokenRepository) DeleteRefreshToken(userID uuid.UUID) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%s", userID.String())

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}

// BlacklistToken adds token to blacklist
func (r *TokenRepository) BlacklistToken(token string, expiration time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", token)

	err := r.client.Set(ctx, key, "blacklisted", expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

// IsTokenBlacklisted checks if token is blacklisted
func (r *TokenRepository) IsTokenBlacklisted(token string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", token)

	_, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}

	return true, nil
}
