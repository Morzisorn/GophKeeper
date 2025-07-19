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

	key, err := getEnvString("SECRET_KEY")
	if err == nil {
		c.SecretKey = key
	}

	cryptoPath, err := getEnvString("CRYPTO_KEY")
	if err == nil {
		c.CryptoKeyPath = cryptoPath
	}
}

func (c *Config) parseAgentEnvs() {

}

func (c *Config) parseServerEnvs() {

	cryptoPath, err := getEnvString("CRYPTO_KEY")
	if err == nil {
		c.CryptoKeyPath = cryptoPath
	}

	db, err := getEnvString("DATABASE_URI")
	if err == nil {
		c.DBConnStr = db
	}

	key, err := getEnvString("SECRET_KEY")
	if err == nil {
		c.SecretKey = key
	}
}

func getEnvString(key string) (string, error) {
	env := os.Getenv(key)
	if env != "" {
		return env, nil
	}
	return "", fmt.Errorf("env %s not found", key)
}

