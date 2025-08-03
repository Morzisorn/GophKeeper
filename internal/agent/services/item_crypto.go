package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gophkeeper/models"
	"io"
)

const (
	nonceSize = 12 // GCM nonce size
)

// EncryptItem encrypts client Item into EncryptedItem
func (cs *CryptoService) encryptItem(item *models.Item) (*models.EncryptedItem, error) {
	mp, err := cs.cnfg.GetMasterPassword()
	if err != nil {
		return nil, err
	}
	if len(mp) == 0 {
		return nil, errors.New("master password not set")
	}

	salt, err := cs.cnfg.GetSalt()
	if err != nil {
		return nil, err
	}
	if len(salt) == 0 {
		return nil, errors.New("user salt not set")
	}

	encryptedData, err := cs.encryptItemData(item.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt item data: %w", err)
	}

	return &models.EncryptedItem{
		ID:            item.ID,
		UserLogin:     item.UserLogin,
		Name:          item.Name,
		Type:          item.Type,
		EncryptedData: *encryptedData,
		Meta:          item.Meta,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}, nil
}

func (cs *CryptoService) decryptItem(encryptedItem *models.EncryptedItem) (*models.Item, error) {
	mp, err := cs.cnfg.GetMasterPassword()
	if err != nil {
		return nil, err
	}
	if len(mp) == 0 {
		return nil, errors.New("master password not set")
	}
	salt, err := cs.cnfg.GetSalt()
	if err != nil {
		return nil, err
	}
	if len(salt) == 0 {
		return nil, errors.New("user salt not set")
	}

	// Create instance of correct data type
	data, err := encryptedItem.Type.CreateDataByType()
	if err != nil {
		return nil, err
	}

	// Decrypt data
	err = cs.decryptData(&encryptedItem.EncryptedData, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt item data: %w", err)
	}

	return &models.Item{
		ID:        encryptedItem.ID,
		UserLogin: encryptedItem.UserLogin,
		Name:      encryptedItem.Name,
		Type:      encryptedItem.Type,
		Data:      data,
		Meta:      encryptedItem.Meta,
		CreatedAt: encryptedItem.CreatedAt,
		UpdatedAt: encryptedItem.UpdatedAt,
	}, nil
}

func (cs *CryptoService) encryptItemData(data models.Data) (*models.EncryptedData, error) {
	// Serialize data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	mk, err := cs.cnfg.GetMasterKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get master key: %w", err)
	}
	if len(mk) == 0 {
		if _, err := cs.generateMasterKey(); err != nil {
			return nil, fmt.Errorf("failed to generate master key: %w", err)
		}
	}

	// Create AES cipher
	block, err := aes.NewCipher(mk)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nil, nonce, jsonData, nil)

	// Return EncryptedData with Base64 encoding
	return &models.EncryptedData{
		EncryptedContent: base64.StdEncoding.EncodeToString(ciphertext),
		Nonce:            base64.StdEncoding.EncodeToString(nonce),
	}, nil
}

func (cs *CryptoService) decryptData(encryptedData *models.EncryptedData, result models.Data) error {
	// Decode Base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData.EncryptedContent)
	if err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(encryptedData.Nonce)
	if err != nil {
		return fmt.Errorf("failed to decode nonce: %w", err)
	}

	// Use cached key
	mk, err := cs.cnfg.GetMasterKey()
	if err != nil {
		return fmt.Errorf("failed to get master key: %w", err)
	}
	if len(mk) == 0 {
		_, err := cs.generateMasterKey()
		if err != nil {
			return fmt.Errorf("failed to generate master key: %w", err)
		}
	}

	// Create AES cipher
	block, err := aes.NewCipher(mk)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Deserialize JSON into result object
	if err := json.Unmarshal(plaintext, result); err != nil {
		return fmt.Errorf("failed to unmarshal decrypted data: %w", err)
	}

	return nil
}
