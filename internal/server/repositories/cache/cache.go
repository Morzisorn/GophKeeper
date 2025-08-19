package cache

import (
	"context"
	"fmt"
	"gophkeeper/config"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisCache struct {
	client     *redis.Client
	defaultTTL time.Duration
	logger     *zap.Logger
}

func NewRedisCache(cfg *config.ServerRedisConfig, logger *zap.Logger) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.Database,
		MinIdleConns:    cfg.MinIdleConns,
		PoolSize:        cfg.PoolSize,
		ConnMaxLifetime: cfg.ConnMaxLifeTime,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Error("redis ping failed", zap.Error(err))
		return nil, err
	}

	logger.Info("redis ping success", zap.String("host", cfg.Host), zap.String("port", cfg.Port))

	return &RedisCache{
		client:     rdb,
		defaultTTL: cfg.DefaultTTL,
		logger:     logger,
	}, nil
}
