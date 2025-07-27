package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"gophkeeper/config"
	"io"
	"os"
	"path/filepath"
)

const (
	privateKeyFilename = "private_key.pem"
	publicKeyFilename  = "public_key.pem"
)

type RSAKeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func LoadRSAKeyPair() error {
	private, err := getKeyPEMFromFile(privateKeyFilename)
	if err != nil {
		return fmt.Errorf("get private key pem from file error: %w", err)
	}

	public, err := getKeyPEMFromFile(publicKeyFilename)
	if err != nil {
		return fmt.Errorf("get public key pem from file error: %w", err)
	}

	cnfg, err := config.GetServerConfig()
	if err != nil {
		return fmt.Errorf("get server config error: %w", err)
	}

	switch {
	// Convert private key from PEM and load private and public keys to config
	case len(private) != 0 && len(public) != 0:
		cnfg.PrivateKey, err = getPrivateKeyFromPEM(private)
		if err != nil {
			return fmt.Errorf("convert private key pem to rsa error: %w", err)
		}
		cnfg.PublicKeyPEM = public
	case len(private) == 0 && len(public) != 0 || len(private) != 0 && len(public) == 0:
		return errors.New("one of public or private keys is destroyed")
	// Generate keys, save in file and load to config
	case len(private) == 0 && len(public) == 0:
		pair, err := generateRSAKeyPair()
		if err != nil {
			return fmt.Errorf("generate rsa keys error: %w", err)
		}

		if err = pair.saveKeysInFile(); err != nil {
			return fmt.Errorf("save rsa keys to file error: %w", err)
		}

		publicPEM, err := pair.getPublicKeyPEM()
		if err != nil {
			return fmt.Errorf("get PEM from public key error: %w", err)
		}
		cnfg.PublicKeyPEM = publicPEM
		cnfg.PrivateKey = pair.PrivateKey
		return nil
	}

	return nil
}

func generateRSAKeyPair() (*RSAKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generate private key error: %w", err)
	}

	return &RSAKeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

func getKeyPEMFromFile(filename string) ([]byte, error) {
	keysPath, err := getKeysPath()
	if err != nil {
		return nil, fmt.Errorf("get keys path error: %w", err)
	}
	filep := filepath.Join(keysPath, filename)

	file, err := os.OpenFile(filep, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file with pem keys: %w. Path: %s", err, filep)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file with pem keys: %w", err)
	}

	return data, nil
}

func (kp *RSAKeyPair) getPublicKeyPEM() ([]byte, error) {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(kp.PublicKey)

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return publicKeyPEM, nil
}

func (kp *RSAKeyPair) getPrivateKeyPEM() []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(kp.PrivateKey)

	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
}

func (kp *RSAKeyPair) saveKeysInFile() error {
	privatePEM := kp.getPrivateKeyPEM()
	if err := saveKeyInFile(privatePEM, privateKeyFilename); err != nil {
		return fmt.Errorf("write private key pem to file error: %w", err)
	}

	publicPEM, err := kp.getPublicKeyPEM()
	if err != nil {
		return fmt.Errorf("get public key PEM error: %w", err)
	}

	if err = saveKeyInFile(publicPEM, publicKeyFilename); err != nil {
		return fmt.Errorf("write public key pem to file error: %w", err)
	}

	return nil
}

func saveKeyInFile(key []byte, filename string) error {
	keysPath, err := getKeysPath()
	if err != nil {
		return fmt.Errorf("get keys path error: %w", err)
	}

	filePath := filepath.Join(keysPath, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file for key pem error: %w", err)
	}
	defer file.Close()

	if _, err = file.Write(key); err != nil {
		return fmt.Errorf("write key pem to file error: %w", err)
	}
	return nil
}

func getKeysPath() (string, error) {
	return config.GetProjectRoot()
}

func GetPublicKeyFromPEM(pemData []byte) (*rsa.PublicKey, error) {
	var zero *rsa.PublicKey
	rest := pemData
	for {
		block, remaining := pem.Decode(rest)
		if block == nil {
			return nil, fmt.Errorf("no PUBLIC KEY block found")
		}

		if _, isPub := any(zero).(*rsa.PublicKey); isPub && block.Type == "PUBLIC KEY" {
			pubkix, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			pub, ok := pubkix.(*rsa.PublicKey)
			if !ok {
				return nil, err
			}
			return pub, nil
		}

		rest = remaining
	}
}

func getPrivateKeyFromPEM(pemData []byte) (*rsa.PrivateKey, error) {
	var zero *rsa.PrivateKey
	rest := pemData
	for {
		block, remaining := pem.Decode(rest)
		if block == nil {
			return nil, fmt.Errorf("no RSA PRIVATE KEY block found")
		}

		if _, isPriv := any(zero).(*rsa.PrivateKey); isPriv && block.Type == "RSA PRIVATE KEY" {
			priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			return priv, nil
		}

		rest = remaining
	}
}
