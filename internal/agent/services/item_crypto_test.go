package services

import (
	"gophkeeper/config"
	"gophkeeper/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptoService_EncryptItem_NilService(t *testing.T) {
	var service *CryptoService = nil

	item := &models.Item{}

	assert.Panics(t, func() {
		service.EncryptItem(item)
	})
}

func TestCryptoService_EncryptItem_EmptyMasterPassword(t *testing.T) {
	service := &CryptoService{
		config: &config.Config{
			AgentConfig: config.AgentConfig{
				MasterPassword: "",
				Salt:           []byte("test-salt"),
			},
		},
	}

	item := &models.Item{}

	result, err := service.EncryptItem(item)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "master password not set")
}

func TestCryptoService_EncryptItem_EmptySalt(t *testing.T) {
	service := &CryptoService{
		config: &config.Config{
			AgentConfig: config.AgentConfig{
				MasterPassword: "test-password",
				Salt:           []byte{},
			},
		},
	}

	item := &models.Item{}

	result, err := service.EncryptItem(item)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user salt not set")
}

func TestCryptoService_DecryptItem_NilService(t *testing.T) {
	var service *CryptoService = nil

	encItem := &models.EncryptedItem{}

	assert.Panics(t, func() {
		service.DecryptItem(encItem)
	})
}

func TestCryptoService_DecryptItem_EmptyMasterPassword(t *testing.T) {
	service := &CryptoService{
		config: &config.Config{
			AgentConfig: config.AgentConfig{
				MasterPassword: "",
				Salt:           []byte("test-salt"),
			},
		},
	}

	encItem := &models.EncryptedItem{}

	result, err := service.DecryptItem(encItem)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "master password not set")
}

func TestCryptoService_encryptData_NilService(t *testing.T) {
	var service *CryptoService = nil

	// Создаем mock данные (нужна реализация models.Data)
	var data models.Data = &MockData{}

	assert.Panics(t, func() {
		service.encryptData(data)
	})
}

func TestCryptoService_decryptData_NilService(t *testing.T) {
	var service *CryptoService = nil

	encData := &models.EncryptedData{}
	var result models.Data = &MockData{}

	assert.Panics(t, func() {
		service.decryptData(encData, result)
	})
}

// Mock структуры для тестирования
type CryptoConfig struct {
	MasterPassword string
	Salt           []byte
	MasterKey      []byte
}

type MockData struct {
	TestField string `json:"test_field"`
}

func (md *MockData) GetType() models.ItemType {
	return models.ItemType(md.TestField)
}

// Реализуем интерфейс models.Data если он существует
// (это зависит от вашей реализации models.Data)
