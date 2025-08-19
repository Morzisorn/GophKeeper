package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func loadEnvFile(envPath string) error {
	return godotenv.Load(envPath)
}

func (c *Config) parseCommonEnvs() {
	addr, err := getEnvString("ADDRESS")
	if err == nil {
		c.Addr = addr
	}
}

func (c *Config) parseAgentEnvs() {}

func (c *Config) parseServerEnvs() {
	db, err := getEnvString("DATABASE_URI")
	if err == nil {
		c.DBConnStr = db
	}
}

func (c *Config) parseServerRedisEnvs() {
	host, err := getEnvString("REDIS_HOST")
	if err == nil {
		c.Redis.Host = host
	}

	port, err := getEnvString("REDIS_PORT")
	if err == nil {
		c.Redis.Port = port
	}

	password, err := getEnvString("REDIS_PASSWORD")
	if err == nil {
		c.Redis.Password = password
	}

	db, err := getEnvInt64("REDIS_URI")
	if err == nil {
		c.Redis.Database = int(db)
	}

	poolSize, err := getEnvInt64("REDIS_POOL_SIZE")
	if err == nil {
		c.Redis.PoolSize = int(poolSize)
	}

	minIdle, err := getEnvInt64("REDIS_MIN_IDLE_CONNS")
	if err == nil {
		c.Redis.MinIdleConns = int(minIdle)
	}

	connsMaxLifeTime, err := getEnvInt64("REDIS_CONNS_MAX_LIFE_TIME")
	if err == nil {
		c.Redis.ConnMaxLifeTime = time.Duration(connsMaxLifeTime)
	}

	connMaxIdleTime, err := getEnvInt64("REDIS_CONN_MAX_IDLE_TIME")
	if err == nil {
		c.Redis.ConnMaxIdleTime = time.Duration(connMaxIdleTime)
	}

	defaultTTL, err := getEnvInt64("REDIS_DEFAULT_TTL")
	if err == nil {
		c.Redis.DefaultTTL = time.Duration(defaultTTL)
	}
}

func getEnvString(key string) (string, error) {
	env := os.Getenv(key)
	if env != "" {
		return env, nil
	}
	return "", fmt.Errorf("env %s not found", key)
}

func getEnvInt64(key string) (int64, error) {
	env := os.Getenv(key)
	if env == "" {
		return 0, fmt.Errorf("env %s not found", key)
	}

	return strconv.ParseInt(env, 10, 64)
}
