package config

import (
	"crypto/rsa"
	"fmt"
	"gophkeeper/internal/logger"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

type Config struct {
	CommonConfig
	AgentConfig
	ServerConfig
}

type CommonConfig struct {
	AppType       string
	Addr          string
	CryptoKeyPath string
	SecretKey     string
}

type AgentConfig struct {
	PublicKey *rsa.PublicKey
}

type ServerConfig struct {
	DBConnStr  string
	PrivateKey *rsa.PrivateKey
	PublicKeyPEM  []byte
}

var (
	instanceServer *Config
	onceServer     sync.Once

	instanceAgent *Config
	onceAgent     sync.Once
)

func GetServerConfig() *Config {
	onceServer.Do(func() {
		var err error
		instanceServer, err = newServerConfig()
		if err != nil {
			logger.Log.Error("Error getting service", zap.Error(err))
		}
	})
	return instanceServer
}

var getEnvPath = getEncFilePath

func newServerConfig() (*Config, error) {
	envPath := getEnvPath()
	if err := loadEnvFile(envPath); err != nil {
		fmt.Printf("Load .env error: %v. Env path: %s\n", err, envPath)
	}

	c := &Config{}

	c.parseCommonEnvs()
	c.parseServerEnvs()

	//c.AppType = "server"

	return c, nil
}

func GetAgentConfig() *Config {
	onceAgent.Do(func() {
		var err error
		instanceAgent, err = newAgentConfig()
		if err != nil {
			logger.Log.Error("Error getting service", zap.Error(err))
		}
	})
	return instanceAgent
}

func newAgentConfig() (*Config, error) {
	envPath := getEnvPath()
	if err := loadEnvFile(envPath); err != nil {
		fmt.Printf("Load .env error: %v. Env path: %s\n", err, envPath)
	}

	c := &Config{}

	c.parseCommonEnvs()
	c.parseAgentEnvs()

	//c.AppType = "agent"

	return c, nil
}

func GetProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd, nil
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			return "", fmt.Errorf("project root not found")
		}
		wd = parent
	}
}

func getEncFilePath() string {
	basePath, err := GetProjectRoot()
	if err != nil {
		logger.Log.Error("Error getting project root ", zap.Error(err))
		return ".env"
	}
	return filepath.Join(basePath, "config", ".env")
}
