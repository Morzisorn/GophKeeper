package config

import (
	"fmt"
	"os"

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

func getEnvString(key string) (string, error) {
	env := os.Getenv(key)
	if env != "" {
		return env, nil
	}
	return "", fmt.Errorf("env %s not found", key)
}
