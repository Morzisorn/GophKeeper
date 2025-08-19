package config

import "time"

type ServerRedisConfig struct {
	Host     string
	Port     string
	Password string
	Database int

	//Connection settings
	PoolSize        int
	MinIdleConns    int
	ConnMaxLifeTime time.Duration
	ConnMaxIdleTime time.Duration

	DefaultTTL time.Duration
}

var defaultServerRedisConfig = ServerRedisConfig{
	Host:            "127.0.0.1",
	Port:            "6379",
	Password:        "redis_password",
	Database:        0,
	PoolSize:        10,
	MinIdleConns:    5,
	ConnMaxLifeTime: time.Hour,
	ConnMaxIdleTime: time.Hour,
	DefaultTTL:      time.Hour,
}

func loadDefaultRedisConfig() *ServerRedisConfig {
	return &defaultServerRedisConfig
}
