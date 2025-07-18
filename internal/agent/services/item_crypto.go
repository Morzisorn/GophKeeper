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
	keySize   = 32 // AES-256
	nonceSize = 12 // GCM nonce size
	saltSize  = 16 // Salt size for PBKDF2
)

// EncryptItem шифрует клиентский Item в EncryptedItem
func (c *CryptoService) EncryptItem(item *models.Item) (*models.EncryptedItem, error) {
	if c.config.MasterPassword == "" {
		return nil, errors.New("master password not set")
	}
	if len(c.config.Salt) == 0 {
		return nil, errors.New("user salt not set")
	}

	encryptedData, err := c.encryptData(item.Data)
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

func (c *CryptoService) DecryptItem(encryptedItem *models.EncryptedItem) (*models.Item, error) {
	if c.config.MasterPassword == "" {
		return nil, errors.New("master password not set")
	}
	if len(c.config.Salt) == 0 {
		return nil, errors.New("user salt not set")
	}

	// Создаем экземпляр правильного типа данных
	data, err := encryptedItem.Type.CreateDataByType()
	if err != nil {
		return nil, err
	}

	// Расшифровываем данные
	err = c.decryptData(&encryptedItem.EncryptedData, data)
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

// ===== ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ =====

// encryptData шифрует объект данных в EncryptedData
func (c *CryptoService) encryptData(data models.Data) (*models.EncryptedData, error) {
	// Сериализуем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// Используем кешированный ключ
	if len(c.config.MasterKey) == 0 {
		c.GenerateMasterKey()
	}

	// Создаем AES шифр
	block, err := aes.NewCipher(c.config.MasterKey)
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
		EncryptedContent:  base64.StdEncoding.EncodeToString(ciphertext),
		Nonce: base64.StdEncoding.EncodeToString(nonce),
	}, nil
}

// decryptData расшифровывает EncryptedData в объект данных
func (c *CryptoService) decryptData(encryptedData *models.EncryptedData, result models.Data) error {
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
	if len(c.GetMasterKey()) == 0 {
		c.GenerateMasterKey()
	}

	// Создаем AES шифр
	block, err := aes.NewCipher(c.GetMasterKey())
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
