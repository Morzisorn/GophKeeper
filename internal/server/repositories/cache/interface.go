package cache

import (
	"context"
	"time"
)

type Cache interface {
	// Basic operations
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// Batch operations
	MGet(ctx context.Context, keys ...string) ([]string, error)
	MSet(ctx context.Context, pairs map[string]interface{}, ttl time.Duration) error
	MDelete(ctx context.Context, keys ...string) error

	// Hash operations
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key string, values map[string]interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)

	// Utility
	TTL(ctx context.Context, key string) (time.Duration, error)
	FlushAll(ctx context.Context) error
	Close() error
	Ping(ctx context.Context) error
}

// Key builders
func UserSessionKey(userID string) string {
	return "session:user:" + userID
}

func UserItemsKey(userID string) string {
	return "items:user:" + userID
}

func UserPublicKeyKey(userID string) string {
	return "pubkey:user:" + userID
}

func ItemKey(itemID string) string {
	return "item:" + itemID
}
