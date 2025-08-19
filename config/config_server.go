package config

import (
	"crypto/rsa"
	"fmt"
)

type ServerCryptoConfig interface {
	SetPrivateKey(*rsa.PrivateKey) error
	SetPublicKeyPEM([]byte) error
}

type DatabaseConfig interface {
	GetConnectionString() string
}

type ServerServicesConfig interface {
	GetPrivateKey() *rsa.PrivateKey
	GetSecretKey() string
}

type ServerControllersConfig interface {
	GetPrivateKey() *rsa.PrivateKey
	GetPublicKeyPEM() []byte
}

type ServerInterceptorsConfig interface {
	GetSecretKey() string
}

type ServerConfig interface {
	ServerInterceptorsConfig
	ServerControllersConfig
	ServerServicesConfig

	GetAddress() string
}

type serverConfig struct {
	Redis *ServerRedisConfig

	DBConnStr    string
	PrivateKey   *rsa.PrivateKey
	PublicKeyPEM []byte
}

func NewServerConfig() (*Config, error) {
	envPath := getEnvPath()
	if err := loadEnvFile(envPath); err != nil {
		fmt.Printf("Load .env error: %v. Env path: %s\n", err, envPath)
	}

	c := &Config{}

	c.Redis = loadDefaultRedisConfig()

	c.parseCommonEnvs()
	c.parseServerEnvs()
	c.parseServerRedisEnvs()

	return c, nil
}
