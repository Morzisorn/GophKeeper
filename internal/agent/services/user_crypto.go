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
	cserv "gophkeeper/internal/server/services/crypto_service"

	"golang.org/x/crypto/pbkdf2"
)

type CryptoService struct {
	Client client.Client
	cnfg   config.AgentCryptoServiceConfig
}

func NewCryptoService(cnfg config.AgentCryptoServiceConfig, client client.Client) (*CryptoService, error) {
	return &CryptoService{
		Client: client,
		cnfg:   cnfg,
	}, nil
}

func (cs *CryptoService) SetPublicKey() error {
	pem, err := cs.Client.GetPublicKeyPEM(context.Background())
	if err != nil {
		return fmt.Errorf("get public key pem from server error: %w", err)
	}

	key, err := cserv.GetPublicKeyFromPEM([]byte(pem))
	if err != nil {
		return fmt.Errorf("get public key from pem error: %w", err)
	}

	if err := cs.cnfg.SetPublicKey(key); err != nil {
		return fmt.Errorf("set public key error: %w", err)
	}

	return nil
}

func (cs *CryptoService) encryptData(data []byte) ([]byte, error) {
	pubkey, err := cs.cnfg.GetPublicKey()
	if err != nil {
		return nil, fmt.Errorf("get public key error: %w", err)
	}
	return rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pubkey,
		data,
		nil,
	)
}

func (cs *CryptoService) setSalt(salt string) error {
	var err error
	saltEnc, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return fmt.Errorf("set salt error: %w", err)
	}

	if err := cs.cnfg.SetSalt(saltEnc); err != nil {
		return fmt.Errorf("set salt error: %w", err)
	}
	mk, err := cs.cnfg.GetMasterKey()
	if err != nil {
		return fmt.Errorf("get master key error: %w", err)
	}
	if len(mk) > 0 {
		mc, err := cs.generateMasterKey()
		if err != nil {
			return fmt.Errorf("generate master key error: %w", err)
		}
		if err := cs.cnfg.SetMasterKey(mc); err != nil {
			return fmt.Errorf("set master key error: %w", err)
		}
	}

	return nil
}

func (cs *CryptoService) setMasterPassword(masterPassword string) error {
	if err := cs.cnfg.SetMasterPassword(masterPassword); err != nil {
		return fmt.Errorf("set master password error: %w", err)
	}
	mk, err := cs.cnfg.GetMasterKey()
	if err != nil {
		return fmt.Errorf("get master key error: %w", err)
	}
	if len(mk) > 0 {
		mc, err := cs.generateMasterKey()
		if err != nil {
			return fmt.Errorf("generate master key error: %w", err)
		}
		if err := cs.cnfg.SetMasterKey(mc); err != nil {
			return fmt.Errorf("set master key error: %w", err)
		}
	}
	return nil
}

func (cs *CryptoService) generateMasterKey() ([]byte, error) {
	mp, err := cs.cnfg.GetMasterPassword()
	if err != nil {
		return nil, fmt.Errorf("get master password error: %w", err)
	}
	salt, err := cs.cnfg.GetSalt()
	if err != nil {
		return nil, fmt.Errorf("get salt error: %w", err)
	}
	mk := pbkdf2.Key([]byte(mp), salt, 10000, 32, sha256.New)
	if err := cs.cnfg.SetMasterKey(mk); err != nil {
		return nil, fmt.Errorf("set master key error: %w", err)
	}
	return mk, nil
}
