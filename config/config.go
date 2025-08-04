package config

import (
	"crypto/rsa"
	"fmt"
	"go.uber.org/zap"
	"gophkeeper/internal/logger"
	"os"
	"path/filepath"
)

type Config struct {
	commonConfig
	agentConfig
	serverConfig
}

type commonConfig struct {
	Addr      string
	SecretKey string
}

func (c *Config) GetConnectionString() string    { return c.DBConnStr }
func (c *Config) GetPrivateKey() *rsa.PrivateKey { return c.PrivateKey }
func (c *Config) GetSecretKey() string           { return c.SecretKey }
func (c *Config) GetPublicKeyPEM() []byte        { return c.PublicKeyPEM }
func (c *Config) GetAddress() string             { return c.Addr }
func (c *Config) SetPrivateKey(pk *rsa.PrivateKey) error {
	if pk == nil {
		return fmt.Errorf("private key is nil")
	}
	c.PrivateKey = pk
	return nil
}
func (c *Config) SetPublicKeyPEM(pk []byte) error {
	if pk == nil {
		return fmt.Errorf("public key is nil")
	}
	c.PublicKeyPEM = pk
	return nil
}

var getEnvPath = getEncFilePath

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
