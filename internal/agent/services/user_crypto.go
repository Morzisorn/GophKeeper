package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"gophkeeper/config"
	"gophkeeper/internal/agent/client"
	"gophkeeper/internal/server/crypto"

	"golang.org/x/crypto/pbkdf2"
)

type CryptoService struct {
	Client client.Client
	config *config.Config
}

func NewCryptoService(client client.Client) (*CryptoService, error) {
	cnfg, err := config.GetAgentConfig()
	if err != nil {
		return nil, fmt.Errorf("get agent config: %w", err)
	}
	return &CryptoService{
		Client: client,
		config: cnfg,
	}, nil
}

func (cs *CryptoService) SetPublicKey() error {
	pem, err := cs.Client.GetPublicKeyPEM(context.Background())
	if err != nil {
		return fmt.Errorf("get public key pem from server error: %w", err)
	}

	key, err := crypto.GetPublicKeyFromPEM([]byte(pem))
	if err != nil {
		return fmt.Errorf("get public key from pem error: %w", err)
	}

	cs.config.PublicKey = key

	return nil
}

func encryptData(data []byte) ([]byte, error) {
	cnfg, err := config.GetAgentConfig()
	if err != nil {
		return nil, fmt.Errorf("get agent config: %w", err)
	}
	return rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		cnfg.PublicKey,
		data,
		nil,
	)
}

func (cs *CryptoService) SetSalt(salt string) error {
	var err error
	cs.config.Salt, err = base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return fmt.Errorf("set salt error: %w", err)
	}

	if len(cs.config.MasterKey) > 0 {
		mc := cs.GenerateMasterKey()
		cs.SetMasterKey(mc)
	}

	return nil
}

func (cs *CryptoService) GetSalt() []byte {
	return cs.config.Salt
}

func (cs *CryptoService) SetMasterPassword(masterPassword string) {
	cs.config.MasterPassword = masterPassword
	if len(cs.config.MasterKey) > 0 {
		mc := cs.GenerateMasterKey()
		cs.SetMasterKey(mc)
	}
}

func (cs *CryptoService) GetMasterPassword() string {
	return cs.config.MasterPassword
}

func (cs *CryptoService) GenerateMasterKey() []byte {
	cs.config.MasterKey = pbkdf2.Key([]byte(cs.config.MasterPassword), cs.config.Salt, 10000, 32, sha256.New)
	return cs.config.MasterKey
}

func (cs *CryptoService) SetMasterKey(mc []byte) {
	cs.config.MasterKey = mc
}

func (cs *CryptoService) GetMasterKey() []byte {
	return cs.config.MasterKey
}
