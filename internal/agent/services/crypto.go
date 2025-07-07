package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"gophkeeper/config"
	"gophkeeper/internal/agent/client"
	"gophkeeper/internal/server/crypto"
)

type CryptoService struct {
	Client client.Client
}

func NewCryptoService(client client.Client) *CryptoService {
	return &CryptoService{
		Client: client,
	}
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

	config.GetAgentConfig().PublicKey = key

	return nil
}

func encryptData(data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		config.GetAgentConfig().PublicKey,
		data,
		nil,
	)
}
