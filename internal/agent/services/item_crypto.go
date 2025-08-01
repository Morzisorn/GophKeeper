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

// EncryptItem шифрует клиентский Item в EncryptedItem
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

	// Создаем экземпляр правильного типа данных
	data, err := encryptedItem.Type.CreateDataByType()
	if err != nil {
		return nil, err
	}

	// Расшифровываем данные
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
	// Сериализуем данные в JSON
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

	// Создаем AES шифр
	block, err := aes.NewCipher(mk)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Создаем GCM режим
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Генерируем случайный nonce
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Шифруем данные
	ciphertext := gcm.Seal(nil, nonce, jsonData, nil)

	// Возвращаем EncryptedData с Base64 кодированием
	return &models.EncryptedData{
		EncryptedContent: base64.StdEncoding.EncodeToString(ciphertext),
		Nonce:            base64.StdEncoding.EncodeToString(nonce),
	}, nil
}

func (cs *CryptoService) decryptData(encryptedData *models.EncryptedData, result models.Data) error {
	// Декодируем Base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData.EncryptedContent)
	if err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(encryptedData.Nonce)
	if err != nil {
		return fmt.Errorf("failed to decode nonce: %w", err)
	}

	// Используем кешированный ключ
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

	// Создаем AES шифр
	block, err := aes.NewCipher(mk)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Создаем GCM режим
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Расшифровываем данные
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Десериализуем JSON в результирующий объект
	if err := json.Unmarshal(plaintext, result); err != nil {
		return fmt.Errorf("failed to unmarshal decrypted data: %w", err)
	}

	return nil
}
