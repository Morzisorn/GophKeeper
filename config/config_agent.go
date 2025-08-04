package config

import (
	"crypto/rsa"
	"fmt"
)

type AgentClientConfig interface {
	GetAddress() string
}

type AgentCryptoServiceConfig interface {
	SetPublicKey(key *rsa.PublicKey) error
	GetPublicKey() (*rsa.PublicKey, error)
	SetSalt(salt []byte) error
	GetSalt() ([]byte, error)
	SetMasterKey(key []byte) error
	GetMasterKey() ([]byte, error)
	SetMasterPassword(masterPassword string) error
	GetMasterPassword() (string, error)
}

func (c *Config) SetPublicKey(key *rsa.PublicKey) error {
	if key == nil {
		return fmt.Errorf("public key is nil")
	}
	c.PublicKey = key
	return nil
}

func (c *Config) GetPublicKey() (*rsa.PublicKey, error) {
	return c.PublicKey, nil
}

func (c *Config) SetSalt(salt []byte) error {
	c.Salt = salt
	return nil
}

func (c *Config) GetSalt() ([]byte, error) {
	if c.Salt == nil {
		return nil, fmt.Errorf("salt is empty")
	}
	return c.Salt, nil
}

func (c *Config) SetMasterKey(key []byte) error {
	c.MasterKey = key
	return nil
}

func (c *Config) GetMasterKey() ([]byte, error) {
	return c.MasterKey, nil
}

func (c *Config) SetMasterPassword(masterPassword string) error {
	c.MasterPassword = masterPassword
	return nil
}

func (c *Config) GetMasterPassword() (string, error) {
	if c.MasterPassword == "" {
		return "", fmt.Errorf("master password is empty")
	}
	return c.MasterPassword, nil
}

type AgentUserServiceConfig interface {
	SetSalt(salt []byte) error
	SetMasterKey(key []byte) error
	SetMasterPassword(masterPassword string) error
}
type agentConfig struct {
	PublicKey      *rsa.PublicKey
	MasterPassword string
	MasterKey      []byte
	Salt           []byte
}

func NewAgentConfig() (*Config, error) {
	envPath := getEnvPath()
	if err := loadEnvFile(envPath); err != nil {
		fmt.Printf("Load .env error: %v. Env path: %s\n", err, envPath)
	}

	c := &Config{}

	c.parseCommonEnvs()
	c.parseAgentEnvs()

	return c, nil
}
